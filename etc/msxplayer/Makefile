AS=tools/sjasmplus
EMU=/Applications/openMSX.app/Contents/MacOS/openmsx
SRCDIR=src
MAINFILE=$(SRCDIR)/main.asm
BUILD=build
ROM=$(BUILD)/main.rom

build: mkdirs
	$(AS) --sym=build/symbols.txt --msg=all --nofakes --raw=$(ROM) $(MAINFILE)

clean:
	rm -rf build

mkdirs:
	mkdir -p $(BUILD)

run:
	$(EMU) $(ROM)