package image

import (
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"sort"
	"time"
)

type imageSort []*cvm.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	// Public images don't have a creation time
	if a[i].CreatedTime == nil || a[j].CreatedTime == nil {
		return false
	}

	itime, _ := time.Parse(time.RFC3339, *a[i].CreatedTime)
	jtime, _ := time.Parse(time.RFC3339, *a[j].CreatedTime)
	return itime.Before(jtime)
}

func mostRecentImage(images []*cvm.Image) *cvm.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))

	return sortedImages[len(sortedImages)-1]
}
