package mieru

import (
	"fmt"

	"github.com/sagernet/sing-box/option"
	mierupb "github.com/enfein/mieru/v3/pkg/appctl/appctlpb"
	"google.golang.org/protobuf/proto"
)

// UpdateUsers hot-swaps mieru credentials without recreating the inbound
// listener or underlay sockets.
func (h *Inbound) UpdateUsers(users []option.MieruUser) error {
	if len(users) == 0 {
		return fmt.Errorf("mieru users is empty")
	}
	pbUsers := make([]*mierupb.User, 0, len(users))
	names := make([]string, 0, len(users))
	for _, user := range users {
		if user.Name == "" || user.Password == "" {
			return fmt.Errorf("mieru username/password must be non-empty")
		}
		pbUsers = append(pbUsers, &mierupb.User{
			Name:     proto.String(user.Name),
			Password: proto.String(user.Password),
		})
		names = append(names, user.Name)
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.server.UpdateUsers(pbUsers); err != nil {
		return err
	}
	h.userNames = names
	return nil
}
