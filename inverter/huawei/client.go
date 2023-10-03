package huawei

type HuaweiClient interface{}

type HuaweiCredential struct {
	Username string
	Password string
}

type huaweiClient struct {
	URL string
}

func NewHuaweiClient(credential *HuaweiCredential) (HuaweiClient, error) {
	client := &huaweiClient{
		URL: URL_VERSION1,
	}

	return client, nil
}
