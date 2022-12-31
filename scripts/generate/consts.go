package main

import (
	"errors"
	"fmt"
	"strings"
)

func generateConsts(d APIDescription) error {
	consts := strings.Builder{}
	consts.WriteString(`
// THIS FILE IS AUTOGENERATED. DO NOT EDIT.
// Regen by running 'go generate' in the repo root.

package gotgbot

`)

	updateConsts, err := generateUpdateTypeConsts(d)
	if err != nil {
		return fmt.Errorf("failed to generate consts for update types: %w", err)
	}
	consts.WriteString(updateConsts)

	consts.WriteString(generateParseModeConsts())

	stickerTypeConsts, err := generateStickerTypeConsts(d)
	if err != nil {
		return fmt.Errorf("failed to generate consts for sticker types: %w", err)
	}
	consts.WriteString(stickerTypeConsts)

	return writeGenToFile(consts, "gen_consts.go")
}

func generateUpdateTypeConsts(d APIDescription) (string, error) {
	updType, ok := d.Types["Update"]
	if !ok {
		return "", errors.New("missing 'Update' type data")
	}
	out := strings.Builder{}
	out.WriteString("// The consts listed below represent all the update types that can be requested from telegram.\n")
	out.WriteString("const (\n")
	for _, f := range updType.Fields {
		if f.Required {
			// All the update types are optional, so skip required values.
			continue
		}
		constName := "UpdateType" + snakeToTitle(f.Name)
		out.WriteString(writeConst(constName, f.Name))
	}
	out.WriteString(")\n\n")
	return out.String(), nil
}

func generateStickerTypeConsts(d APIDescription) (string, error) {
	updType, ok := d.Types["Sticker"]
	if !ok {
		return "", errors.New("missing 'Sticker' type data")
	}
	out := strings.Builder{}
	out.WriteString("// The consts listed below represent all the sticker types that can be obtained from telegram.\n")
	out.WriteString("const (\n")
	for _, f := range updType.Fields {
		if f.Name != "type" {
			// the field we want to look at is called "type", ignore all others.
			continue
		}
		types, err := extractQuotedValues(f.Description)
		if err != nil {
			return "", fmt.Errorf("failed to get quoted types: %w", err)
		}
		for _, t := range types {
			constName := "StickerType" + snakeToTitle(t)
			out.WriteString(writeConst(constName, t))
		}
	}
	out.WriteString(")\n\n")
	return out.String(), nil
}

func generateParseModeConsts() string {
	// Adding these manually because they're not part of the spec, and theyre not going to change much anyway.
	formattingTypes := []string{"HTML", "MarkdownV2", "Markdown", "None"}

	out := strings.Builder{}
	out.WriteString("// The consts listed below represent all the parse_mode options that can be sent to telegram.\n")
	out.WriteString("const (\n")
	for _, t := range formattingTypes {
		constName := "ParseMode" + t
		if t == "None" {
			// no parsemode == empty string value.
			out.WriteString(writeConst(constName, ""))
			continue
		}
		out.WriteString(writeConst(constName, t))
	}
	out.WriteString(")\n\n")
	return out.String()
}

func writeConst(name string, value string) string {
	return fmt.Sprintf("%s = \"%s\"\n", name, value)
}