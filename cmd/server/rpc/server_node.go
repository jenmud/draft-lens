package main

import (
	"context"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

func (s *server) AddNode(ctx context.Context, req *pb.NodeReq) (*pb.NodeResp, error) {
	kvs := make([]graph.KV, len(req.Properties))

	count := 0
	for k, v := range req.Properties {
		kvs[count] = graph.KV{Key: k, Value: v}
		count++
	}

	node, err := s.graph.AddNode(req.Uid, req.Label, kvs...)
	if err != nil {
		return nil, err
	}

	resp := pb.NodeResp{
		Uid:        node.UID,
		Label:      node.Label,
		Properties: node.Properties,
		InEdges:    node.InEdges(),
		OutEdges:   node.OutEdges(),
	}

	return &resp, nil
}

func (s *server) RemoveNode(ctx context.Context, req *pb.UIDReq) (*pb.RemoveResp, error) {
	var errmsg string

	err := s.graph.RemoveNode(req.Uid)
	if err != nil {
		errmsg = err.Error()
	}

	return &pb.RemoveResp{Uid: req.Uid, Success: err == nil, Error: errmsg}, nil
}

func (s *server) Node(ctx context.Context, req *pb.UIDReq) (*pb.NodeResp, error) {
	node, err := s.graph.Node(req.Uid)
	if err != nil {
		return nil, err
	}

	resp := &pb.NodeResp{
		Uid:        node.UID,
		Label:      node.Label,
		Properties: node.Properties,
		InEdges:    node.InEdges(),
		OutEdges:   node.OutEdges(),
	}

	return resp, nil
}

func (s *server) Nodes(req *pb.NodesReq, stream pb.Graph_NodesServer) error {
	iter := s.graph.Nodes()
	for iter.Next() {
		node := iter.Value().(graph.Node)

		resp := pb.NodeResp{
			Uid:        node.UID,
			Label:      node.Label,
			Properties: node.Properties,
			InEdges:    node.InEdges(),
			OutEdges:   node.OutEdges(),
		}

		if err := stream.Send(&resp); err != nil {
			return nil
		}
	}

	return nil
}
