package config

import (
	"path"
	"strings"

	"gitee.com/golden-go/golden-go/pkg/utils/base_dir"
	"gitee.com/golden-go/golden-go/pkg/utils/ldap"
	"gitee.com/golden-go/golden-go/pkg/utils/logger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	// 16为密码加密
	viper.SetDefault("goldengo.password.key", "KY9ciRr1Q7sOgjVV")
	// mysql连接url
	viper.SetDefault("mysql.dsn", "golden_go:golden_go123@tcp(127.0.0.1:3306)/golden_go?charset=utf8&parseTime=True&loc=Local")
	//监听地址
	viper.SetDefault("listen", ":8080")
	//jwt token失效时间 单位分钟
	viper.SetDefault("jwt.exp", 60)
	//默认公钥
	viper.SetDefault("jwt.publicKey", `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsTlzGXqZPhXiVaDnq4ks
DuWvqyA1/nUauAGlLKDYxIo9WfBvugwrsZR9Sy/b0B7MnyBWlepOPo+8I7OXpKpR
JQ6QiwVO8qAD1uyL/7zyh9dN3QdAcOVHg3KNM3vWUYV/T5Um2hMO69FBbTy8TNeA
iov5Hi9LeJl0Zrz533xQ9VVMPi/LgfM4Fkng1uQtqti2bvTFXwTOufMVBjWMrJcJ
euGTrxqbFlvhrf9947m1SAJm3DNaUS/4cz4eg6aOAprmlHN+H1ND8yhk4nTPHuH6
Jyn0R63NiQTKSM4ZKpUMKv6kua4oOShxywdzk5I6n1tHIuDbfw6jTHmVMKc3m/zZ
ywIDAQAB
-----END PUBLIC KEY-----`)
	//默认私钥
	viper.SetDefault("jwt.privateKey", `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCxOXMZepk+FeJV
oOeriSwO5a+rIDX+dRq4AaUsoNjEij1Z8G+6DCuxlH1LL9vQHsyfIFaV6k4+j7wj
s5ekqlElDpCLBU7yoAPW7Iv/vPKH103dB0Bw5UeDco0ze9ZRhX9PlSbaEw7r0UFt
PLxM14CKi/keL0t4mXRmvPnffFD1VUw+L8uB8zgWSeDW5C2q2LZu9MVfBM658xUG
NYyslwl64ZOvGpsWW+Gt/33jubVIAmbcM1pRL/hzPh6Dpo4CmuaUc34fU0PzKGTi
dM8e4fonKfRHrc2JBMpIzhkqlQwq/qS5rig5KHHLB3OTkjqfW0ci4Nt/DqNMeZUw
pzeb/NnLAgMBAAECggEAdsPoNWfqcCfcQMQO3O6VHvqfqc9xP7AckrlPhsPX4IY/
rRkq1oQ3d87p1EwjQ6qQOAdE3zxg6R4L1+UPt6MHtAy5ouCQ0pmXWR22iFCIO652
mKu2bLHKJfXLWHgh3QnYkm2C1tu0wSW/ccQk8F3u32oELU1Gh3BXCE6WKUW3P2Cn
mIkkPQb96E9yDuFHd0ccx1klNLpmxr/xS5KRdNb5iQyeGT/li8ZPfQ74oIgoycYa
X5NFagPxsnha53R/622xJpMCjzKFZG3bAi2LHEbAaidFHe6IwVKy6yebeicrJhMs
plZqEedQTRIK/gX8qxKXcApHDB/PikIzQAY8snl0uQKBgQDbJtE1oyWGP4dFskr4
srKVZpwUQRdTzAMCEPRg62pXGUhrfcJJmrHWDpP3qbBclMt3tNrlRczfbgq5Rk1c
JxIkLodjRxEkeUSOMa/BbGoWV84Ylcy4Wdw4nAasCGhlwJUDC2YsKWRs+G6CDAEm
yroq8vXqHI7tFIERZ/+ntdKbrQKBgQDPBenVaQpnNyUfyKa3YFfcaXZgEovL3O4/
1PIZ/1HoCUcXRDCc6T4oEk4hhgMFLwzrJTDfhQFZOqzGSUua+gETjvB8x+tcqxii
Zi4E1RlD3vttzRovQf0x17YyvCooy50aOzrQkfm8pZWmV49WH1lBDqRjpYa/zon1
7ho650/6VwKBgDfftLVNLDMHHXEKnQy9WsS2jZKoac8mk5nCQtw3CTy9qHYncRKd
Czw9KUKak+l20k1p4elUm4BXqQpFv1GAcKKi7kfEhb1b4buzNVFGV+HsbxJblv9l
gb05IoZX+m3+5L8K9/jIcN9Lk7k6YEzIoSB0I3iV4WbWtHWAo3DZ7aFdAoGALScn
Byrv3+9BI5uJ8TkMMMC31uB0qTQ4qqdxXoY3gRp07PgseizNJ8RGUj1+byUB9k+Y
R0glMScBBAZ8fQLGNVPZ0/1usjmHF/SdWOR5rDn4MMypR2FnLfXUgWBU+azfSRde
DpuoEqMy7qLuWmRf/TaKWEmECFWt9XvHMM2+veUCgYEAhbXHE7VHu6KDu0JEEM96
6fsDeoC2a8DfxT9FW6PWzmK0jew1VWJ78rV3zlX9pX2JUjNgabcZk5NYtJVsBwNC
e2gKBcM0f+h06uDwbveyUlxgnd6p4i0D69Kg6IiPcgZJWSZ0BvoZ7+MI/ZGLPdx5
qdxS6V5MFi8tWrhRHCo0jGA=
-----END PRIVATE KEY-----
`)
	viper.SetDefault("auth.ldap.enable", false)
	viper.SetDefault("auth.ldap.servers", []*ldap.ServerConfig{})
}

func InitConfig(cfgFile, configNmae string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		logger.Debug("use default cfgfile")
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(home)
		bd := base_dir.GetBaseDir()
		viper.AddConfigPath(bd)
		viper.AddConfigPath(path.Join(bd, "conf"))
		viper.AddConfigPath(path.Join(bd, "etc"))
		if strings.TrimSpace(configNmae) == "" {
			viper.SetConfigName("base-service")
		} else {
			viper.SetConfigName(configNmae)
		}

	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		logger.Info("Using config file:" + viper.ConfigFileUsed())
	} else {
		logger.Warn("read in config", zap.Error(err))
	}
	return nil
}
