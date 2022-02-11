ANTLR_CMD ?= antlr
ANTLR_GENERATED_DIR ?= generated
ANTLR_GRAMMAR_FILE ?= grammar/m4l.g4
ANTLR_OUT_LANG ?= Go

clean:
	rm -rf $(ANTLR_GENERATED_DIR)

grammar:
	$(ANTLR_CMD) -o $(ANTLR_GENERATED_DIR) -Dlanguage=$(ANTLR_OUT_LANG) $(ANTLR_GRAMMAR_FILE)

.PHONY: grammar clean