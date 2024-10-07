package doc

//KVMessage : { "message" : "val" //tag }
func KVMessage(val string, tag ...string) keyValue {
	return KV("message", `"`+val+`"`, tag...)
}

//KV :
func KV(k, v interface{}, tag ...string) keyValue {
	kv := keyValue{
		Key: k,
		Val: v,
	}
	for _, c := range tag {
		kv.Tag += c
		if len(tag) > 1 {
			kv.Tag += "\n"
		}
	} //for
	return kv
}

//KVMessageString : { "message" : string //tag }
func KVMessageString(tag ...string) keyValue {
	return KV("message", String, tag...)
}

//D : { notKey , "value" , "tag"...}
func D(val interface{}, tag ...string) keyValue {
	kv := keyValue{
		Key: notKey,
		Val: val,
	}
	for _, c := range tag {
		kv.Tag += c
		if len(tag) > 1 {
			kv.Tag += "\n"
		}
	} //for
	return kv
}
