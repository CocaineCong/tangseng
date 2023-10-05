package job

import (
	"context"

	"github.com/robfig/cron/v3"
)

func RegisterJob(ctx context.Context) {
	c := cron.New()
	// 至于为什么要把存储颗粒定义这么大，没有周纬度存储，是因为我觉得这个合并是离线任务，颗粒大，意味着最终的库少，意味着查询的时候，会减少对磁盘的io
	// 优化搜索引擎查询的速度，毕竟搜索引擎实时性只有查询。

	// 每周的周日凌晨3点，把这一周的所有天数都合并到这个月中
	_, _ = c.AddJob("0 0 3 ? * SUN", &Command{Name: "MergeInvertedIndexDay2Month", Exec: MergeInvertedIndexDay2Month, Context: ctx})
	// 每个月的最后一天的凌晨5点20分，把这一个月的索引数据都合并到这个季度中
	_, _ = c.AddJob("0 20 5 L * ?", &Command{Name: "MergeInvertedIndexMonth2Season", Exec: MergeInvertedIndexMonth2Season, Context: ctx})

	c.Start()
}
