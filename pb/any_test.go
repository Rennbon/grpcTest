package pb

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/magiconair/properties/assert"
	"grpcTest/pb/entity"
	"testing"
)

func TestProtoAny(t *testing.T) {

	pm := &ClientComMeta{
		Resource: "user",
		Revision: "v1.1",
		Action:   "GetUserInfo",
	}
	object, err := ptypes.MarshalAny(pm)
	if err != nil {
		t.Error(err)
		return
	}

	/*	pipeR := PipeRequest{
			Qid:    "1",
			Params: object,
		}
	*/
	fmt.Println(object)
	pmu := &ClientComMeta{}
	err = ptypes.UnmarshalAny(object, pmu)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, pmu.String(), pm.String())
}
func TestEntityMap(t *testing.T) {
	pm := &entity.UserInfo{
		UserId:   "id123",
		UserName: "Rennbon",
		Face:     "https://upload.jianshu.io/users/upload_avatars/7215015/1d634670-8438-4a49-9748-9fc36f5168ce?imageMogr2/auto-orient/strip|imageView2/1/w/120/h/120",
	}
	object, err := ptypes.MarshalAny(pm)
	if err != nil {
		t.Error(err)
		return
	}
	pmu := &entity.UserInfo{}
	fmt.Println(pm)
	err = ptypes.UnmarshalAny(object, pmu)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pmu.String())
}
