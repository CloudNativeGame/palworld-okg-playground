package gameserver

import (
	gamekruisev1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	"github.com/openkruise/kruise-game/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"testing"
	"time"
)

func TestSortGs(t *testing.T) {
	tests := []struct {
		before []gamekruisev1alpha1.GameServer
		after  []int
	}{
		{
			before: []gamekruisev1alpha1.GameServer{
				{
					ObjectMeta: metav1.ObjectMeta{
						CreationTimestamp: metav1.Time{Time: time.Now().Add(10 * time.Minute)},
						Name:              "xxx-0",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						CreationTimestamp: metav1.Time{Time: time.Now()},
						Name:              "xxx-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						CreationTimestamp: metav1.Time{Time: time.Now().Add(5 * time.Minute)},
						Name:              "xxx-2",
					},
				},
			},
			after: []int{1, 2, 0},
		},
	}

	for caseNum, test := range tests {
		after := SortGs(test.before)
		sort.Sort(after)
		expect := test.after
		actual := util.GetIndexListFromGsList(after)
		for i := 0; i < len(actual); i++ {
			if expect[i] != actual[i] {
				t.Errorf("case %d: expect %v but got %v", caseNum, expect, actual)
			}
		}
	}

}
