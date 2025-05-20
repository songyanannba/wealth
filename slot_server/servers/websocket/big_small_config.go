package websocket

type BigOrSmallConfig struct {
	Seat       int
	BigOrSmall int // 1=大（粉色） 2=小（紫色）
}

func GetBigOrSmallConfigSort() []*BigOrSmallConfig {
	var (
		bigOrSmallConfig = make([]*BigOrSmallConfig, 18)
	)
	//根据微信发的手动排序
	bigOrSmallConfigSort := []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2}
	for i, v := range bigOrSmallConfigSort {
		bigOrSmallConfig[i] = &BigOrSmallConfig{
			Seat:       i,
			BigOrSmall: v,
		}
	}
	return bigOrSmallConfig
}
