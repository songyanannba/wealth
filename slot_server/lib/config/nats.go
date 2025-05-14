package config

const SlotServer = "slot_server"

const (
	AnimalParty = "animal_party" // 流名称

	//主题+消费者 接受网关发送过来的消息

	AnimalPartyTopic        = "animal.party.topic1"   // 流绑定的主题
	AnimalPartyConsumerName = "animal_party_consumer" //消费者

	//主题+消费者 处理完成 发送到网关 网关消费者接收并返回客户端

	AnimalPartyTopicResp           = "animal.party.topic.resp1"   // 流绑定的 主题
	AnimalPartyProducerSubjectResp = "animal_party_resp_consumer" // 消费者
)
