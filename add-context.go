package main

func ProcessRawCaptions(url string, rawCaptions []RawCaption) []Caption {
	captions := make([]Caption, len(rawCaptions))

	// Iterate over the rawCaptions
	for i, rawCaption := range rawCaptions {
		var context string

		// Add previous RawCaption text to the context, if available
		if i > 0 {
			context += rawCaptions[i-1].Text + " "
		}

		// Add current RawCaption text to the context
		context += rawCaption.Text + " "

		// Add next RawCaption text to the context, if available
		if i < len(rawCaptions)-1 {
			context += rawCaptions[i+1].Text
		}

		caption := Caption{
			Url:           url,
			Index:         rawCaption.Index,
			Text:          rawCaption.Text,
			Context:       context,
			TimestampFrom: rawCaption.TimestampFrom,
			TimestampTo:   rawCaption.TimestampTo,
			ClipLength:    rawCaption.ClipLength,
		}

		captions[i] = caption
	}

	return captions
}
