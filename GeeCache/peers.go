package geecache

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 建立的注册节点？借助一致性哈希进行选择节点；实现http客户端，建立与远程节点的通信。

// 建立两个接口，用于实现节点选择和客户端建立过程中的类型设计

// 用于选取key对应的节点PeerGetter

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 用于从对应group缓存中查找缓存只，

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

// 实现客户端功能

type httpGetter struct {
	baseUrl string
}

// url标注库的第一次使用？

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	// 从http服务中获取u对应的内容
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 这里有一个小的注意点，这里一定是要进行resp的err是否为空这一内容进行判断才有必要进行defer关闭的，否则没有获取到数据就是对一个空的接口进行关闭，执行时会发生错误

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	// 读取获取到的数据
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	// 将获取到的数据进行回传
	return bytes, nil
}

// 初始化一个的目的是什么：？

var _ PeerGetter = (*httpGetter)(nil)

// 增加一个选择节点的操作

const (
	//defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		// 说明完成了所需功能
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// 这个的意思是为了检查实现了当前接口的所有方法

var _ PeerPicker = (*HTTPPool)(nil)

// 实现的功能：访问远端节点信息，还要集成到主程序服务中去
