NAME ?= twittersounds
VERSION = `git describe`

zip:
		zip -r $(NAME)-$(VERSION).zip . -x .git/\*

.PHONY: zip
