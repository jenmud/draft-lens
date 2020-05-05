package main

import (
	"context"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

func (s *server) AddEdge(ctx context.Context, req *pb.EdgeReq) (*pb.EdgeResp, error) {
	kvs := make([]graph.KV, len(req.Properties))

	count := 0
	for k, v := range req.Properties {
		kvs[count] = graph.KV{Key: k, Value: v}
		count++
	}

	edge, err := s.graph.AddEdge(req.Uid, req.SourceUid, req.Label, req.TargetUid, kvs...)
	if err != nil {
		return nil, err
	}

	resp := pb.EdgeResp{
		Uid:        edge.UID,
		SourceUid:  edge.SourceUID,
		Label:      edge.Label,
		TargetUid:  edge.TargetUID,
		Properties: edge.Properties,
	}

	return &resp, nil
}

func (s *server) RemoveEdge(ctx context.Context, req *pb.UIDReq) (*pb.RemoveResp, error) {
	var errmsg string

	err := s.graph.RemoveEdge(req.Uid)
	if err != nil {
		errmsg = err.Error()
	}

	return &pb.RemoveResp{Uid: req.Uid, Success: err == nil, Error: errmsg}, nil
}

func (s *server) Edge(ctx context.Context, req *pb.UIDReq) (*pb.EdgeResp, error) {
	edge, err := s.graph.Edge(req.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.EdgeResp{Uid: edge.UID, SourceUid: edge.SourceUID, Label: edge.Label, TargetUid: edge.TargetUID, Properties: edge.Properties}, nil
}

func (s *server) Edges(req *pb.EdgesReq, stream pb.Graph_EdgesServer) error {
	iter := s.graph.Edges()
	for iter.Next() {
		edge := iter.Value().(graph.Edge)
		resp := pb.EdgeResp{Uid: edge.UID, SourceUid: edge.SourceUID, Label: edge.Label, TargetUid: edge.TargetUID, Properties: edge.Properties}
		if err := stream.Send(&resp); err != nil {
			return nil
		}
	}

	return nil
}
