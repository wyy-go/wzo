package main

import (
	"strings"

	"github.com/wyy-go/wlib/proc"
	"google.golang.org/protobuf/compiler/protogen"
)

// annotation const value
const (
	Identity             = "errno"
	AttributeNameStatus  = "status"
	AttributeNameCode    = "code"
	AttributeNameMessage = "message"
)

type ErrnoDerive struct {
	Enabled bool
	Status  int
}

func ParseDeriveErrno(s protogen.Comments) (*ErrnoDerive, proc.CommentLines) {
	derives, remainComments := proc.NewCommentLines(s.String()).FindDerives(Identity)
	ds := proc.Derives(derives)
	ret := &ErrnoDerive{
		Enabled: ds.ContainHeadless(Identity),
		Status:  500,
	}
	values := ds.FindValue(Identity, AttributeNameStatus)
	for _, value := range values {
		if v, ok := value.(proc.Integer); ok && v.Value > 0 && v.Value < 1000 {
			ret.Status = int(v.Value)
		}
	}
	return ret, remainComments
}

type ErrnoValueDerive struct {
	Status  int
	Code    int
	Message string
}

func ParseDeriveErrnoValue(status, code int, s protogen.Comments) (*ErrnoValueDerive, proc.CommentLines) {
	derives, remainComments := proc.NewCommentLines(s.String()).FindDerives(Identity)
	ret := &ErrnoValueDerive{
		Status:  status,
		Code:    code,
		Message: strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(remainComments.LineString()), "\n", ","), `"`, `\"`),
	}
	for _, d := range derives {
		for _, v := range d.Attrs {
			switch v.Name {
			case AttributeNameStatus:
				if v, ok := v.Value.(proc.Integer); ok {
					ret.Status = int(v.Value)
				}
			case AttributeNameCode:
				if v, ok := v.Value.(proc.Integer); ok {
					ret.Code = int(v.Value)
				}
			case AttributeNameMessage:
				if v, ok := v.Value.(proc.String); ok {
					ret.Message = v.Value
				}
			}
		}
	}
	return ret, remainComments
}
