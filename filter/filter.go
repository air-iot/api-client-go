package filter

import (
	"context"

	"github.com/go-kratos/kratos/v2/selector"

	"github.com/air-iot/logger"
)

// Metadata is metadata filter.
func Metadata(metadata map[string]string) selector.NodeFilter {
	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		logger.Infof("grpc client context: %+v metadata %+v", ctx, metadata)
		for _, n := range nodes {
			logger.Infof("grpc client metadata nodes, Scheme: %s, Metadata: %+v, Address: %s", n.Scheme(), n.Metadata(), n.Address())
		}
		if len(metadata) == 0 {
			return nodes
		}
		newNodes := nodes[:0]
		for _, n := range nodes {
			f := true
			if len(n.Metadata()) == 0 {
				f = false
			}
			for k, v := range n.Metadata() {
				nv, ok := metadata[k]
				if !ok {
					f = false
					break
				}
				if nv != v {
					f = false
					break
				}
			}
			if f {
				newNodes = append(newNodes, n)
			}
		}
		return newNodes
	}
}
