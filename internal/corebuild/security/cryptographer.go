package security

type Cryptographer interface {
	DecryptYamlContent(cipherTextYaml string) (string, error)
}
