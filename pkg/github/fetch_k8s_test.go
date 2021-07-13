package github

import "testing"

func Test_fetchTags(t *testing.T) {
	tags := fetchTags()
	println(tags)
}

func Test_fetchSealyunTags(t *testing.T) {
	tags := fetchSealyunTags()
	println(tags)
}
