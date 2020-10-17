package protocol

import (
	"errors"
	"strings"
)

const STRING_SEP = "\""
const DATA_TOKENIZER = " "

type Verb struct {
	Type string
	Payload []string
}

func Parse(data string) (verb Verb, err error) {
	tokens := tokenize(data)

	if len(tokens) > 0 {
		switch firstToken := strings.ToUpper(tokens[0]); firstToken {
		case "KIKOO":
			verb, err = createInit(tokens[1:])

		case "ASV":
			verb, err = createAsk()

		case "OKLM":
			verb, err = createIdent(tokens[1:])

		case "TROPA":
			verb, err = createReject()

		case "TAVU":
			verb, err = createMsg(tokens[1:])

		case "LOL":
			verb, err = createAck()

		case "JPP":
			verb, err = createQuit()

		case "WTF":
			verb, err = createErr(tokens[1:])

		default:
			err = errors.New("Unknown verb")
		}
	} else {
		err = errors.New("Empty string")
	}

	return
}

func createInit(tokens []string) (Verb, error) {
	var name string
	name, tokens = popNextString(tokens)
	authority := tokens[0]
	verb := Verb { 
		Type: "INIT", 
		Payload: []string{ name, authority },
	}

	return verb, nil
}

func createAsk() (Verb, error) {
	verb := Verb {
		Type: "ASK",
	}

	return verb, nil
}

func createIdent(tokens []string) (verb Verb, err error) {
	var payload []string

	if len(tokens) == 0 {
		err = errors.New("No name provided")

		return
	}

	var name string
	name, tokens = popNextString(tokens)
	if len(name) == 0 {
		err = errors.New("Empty name")

		return
	}

	payload = append(payload, name)

	if len(tokens) == 0 {
		err = errors.New("No public authority provided")

		return
	}
	
	payload = append(payload, tokens[0])

	if len(tokens) > 1 {
		for _, other := range tokens[1:] {
			if other != "/" {
				payload = append(payload, other)
			}
		}
	}

	verb = Verb {
		Type: "IDENT",
		Payload: payload,
	}

	return
}

func createReject() (Verb, error) {
	verb := Verb {
		Type: "REJECT",
	}

	return verb, nil
}

func createMsg(tokens []string) (Verb, error) {
	var message string
	message, tokens = popNextString(tokens)

	verb := Verb {
		Type: "MSG",
		Payload: []string{ message },
	}

	return verb, nil
}

func createAck() (Verb, error) {
	verb := Verb {
		Type: "ACK",
	}

	return verb, nil
}

func createQuit() (Verb, error) {
	verb := Verb {
		Type: "QUIT",
	}

	return verb, nil
}

func createErr(tokens []string) (Verb, error) {
	var explanation string
	explanation, tokens = popNextString(tokens)

	verb := Verb {
		Type: "ERR",
		Payload: []string{ explanation },
	}

	return verb, nil
}

func tokenize(data string) []string {
	return strings.Split(data, DATA_TOKENIZER)
}

func popNextString(tokens []string) (string, []string) {
	var output []string
	count := 0
	for _, token := range tokens {
		if strings.HasPrefix(token, STRING_SEP) || count > 0 {
			if token == STRING_SEP + STRING_SEP {
				token = STRING_SEP
			}

			output = append(output, token)
			count++
		}

		if strings.HasSuffix(token, STRING_SEP) {
			break;
		}
	}

	return strings.Trim(strings.Join(output, DATA_TOKENIZER), STRING_SEP), tokens[count:]
}