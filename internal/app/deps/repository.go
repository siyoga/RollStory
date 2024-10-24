package deps

import (
	"github.com/siyoga/rollstory/internal/repository"
	"github.com/siyoga/rollstory/internal/repository/story"
	"github.com/siyoga/rollstory/internal/repository/thread"
)

func (d *dependencies) ThreadRepository() repository.ThreadRepository {
	if d.threadRepository == nil {
		d.threadRepository = thread.NewThreadRepository(d.log, d.RedisThreadClient())
	}

	return d.threadRepository
}

func (d *dependencies) StoryRepository() repository.StoryRepository {
	if d.storyRepository == nil {
		d.storyRepository = story.NewThreadRepository(d.log, d.RedisStoryClient())
	}

	return d.storyRepository
}
