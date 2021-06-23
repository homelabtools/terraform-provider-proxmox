### Automatic logging
# Appending $(LOGGED) to a line of a make target will automatically log its output
# to $(BUILD_LOG_DIR) named in the pattern of 'target_<nameOfTarget>_<dateAndTime>'
BUILD_LOG_DIR_NAME := .buildlogs
BUILD_LOG_DIR      := $(shell mkdir -p $(BUILD_LOG_DIR_NAME) && printf $(BUILD_LOG_DIR_NAME))
LOGGED              = 2>&1 | tee -a $(BUILD_LOG_DIR)/target_$(shell basename $@)_$(shell date | tr ' ' '_')