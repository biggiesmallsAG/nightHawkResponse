import re

class ValidateUserInput:
	def __init__(self, user_input):
		self.input = user_input

	def ValidateInputMixed(self): ## Matching only Alpha Numeric on user input.
		return re.match("^[a-zA-Z0-9.\s]+$", self.input)

	def ValidateInputMixedPunctual(self): ## Matching only Alpha Numeric on user input + Punctuation
		return re.match("^[a-zA-Z0-9.\-_]+$", self.input)

	def ValidateInputInteger(self):
		return re.match("^[0-9]+$", self.input)

	def ValidateIPAddr(self):
		return re.match("^(\d{1,3}\.){3}\d{1,3}(\/\d{1,2}|)$", self.input)