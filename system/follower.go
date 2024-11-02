package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/component"
)

type Follower struct {
	query *donburi.Query
}

func NewFollower() *Follower {
	return &Follower{
		query: donburi.NewQuery(filter.Contains(transform.Transform, component.Follower)),
	}
}

func (s *Follower) Update(w donburi.World) {
	s.query.Each(w, func(entry *donburi.Entry) {
		follower := component.Follower.Get(entry)
		if follower.Target == nil || !follower.Target.Valid() {
			return
		}

		follower.FollowingTimer.Update()
		if follower.FollowingTimer.IsReady() {
			follower.Target = nil
			return
		}

		// TODO: Should rather rotate towards the target instead of looking at it straight away.
		targetPos := transform.WorldPosition(follower.Target)
		transform.LookAt(entry, targetPos)
		component.Velocity.Get(entry).Velocity = transform.Right(entry).MulScalar(follower.FollowingSpeed)
	})
}
