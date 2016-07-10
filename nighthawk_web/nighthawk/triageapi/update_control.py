def TagIntToStr(tag):
	TAG_CHOICES = {
		"0": "Benign",
		"1": "Follow Up",
		"2": "Malicious",
		"3": "For Review"
	}

	return TAG_CHOICES[tag]

def FpIntToStr(fp):
	LENGTH_CHOICES = {
		"0": True,
		"1": False
	}

	return LENGTH_CHOICES[fp]