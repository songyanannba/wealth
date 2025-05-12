package config

// NatsMemeBattle meme 服务ID
const NatsMemeBattle = "meme_battle"

const (
	//MemeBattle 流绑定了2个主题 ，每个主题有一个消费者
	MemeBattle = "meme_battle" // 流名称

	//主题+消费者 接受网关发送过来的消息

	MemeBattleTopic = "meme.battle.topic"    // 流绑定的主题
	ConsumerName    = "meme_battle_consumer" //消费者

	//主题+消费者 处理完成 发送到网关 网关消费者接收并返回客户端

	MemeBattleTopicResp = "meme.battle.topic.resp"    // 流绑定的 主题
	ProducerSubjectResp = "meme_battle_resp_consumer" // 消费者
)
