package ipconfig

import "github.com/coderc/im/ipconfig/domain"

func top5EdnPoints(eds []*domain.EndPoint) []*domain.EndPoint {
	if len(eds) < 5 {
		return eds
	}
	return eds[:5]
}

func packRes(ed []*domain.EndPoint) Response {
	return Response{
		Message: "ok",
		Code:    0,
		Data:    ed,
	}
}
