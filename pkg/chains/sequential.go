package chains

type SimpleSecuentialChain struct {
	chains []Chain
}

func NewSimpleSecuentialChain(chains []Chain) SimpleSecuentialChain {
	return SimpleSecuentialChain{chains: chains}
}

func (c *SimpleSecuentialChain) Run(inputVariable string) (string, error) {
	if len(c.chains) == 0 {
		return "", nil
	} else if len(c.chains) == 1 {
		return c.chains[0].Run(inputVariable, nil)
	}

	output, err := c.chains[0].Run(inputVariable, nil)
	if err != nil {
		return "", err
	}
	for _, chain := range c.chains[0:] {
		resp, err := chain.Run(output, nil)
		if err != nil {
			return "", err
		}
		output = resp
	}

	return output, nil
}
