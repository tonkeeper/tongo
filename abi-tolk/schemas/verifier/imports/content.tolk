// StringTail matches FunC's storeStringTail/loadStringTail - snake-encoded string stored
// inline in the current cell (not like ordinary Tolk's string which is basically Cell<StringTail>)
type StringTail = string

fun StringTail.unpackFromSlice(mutate s: slice): StringTail {
    return "";
}

fun StringTail.packToBuilder(self, mutate b: builder) {
}

// SourceContent is a verified source payload stored as a separate cell ref.
// From contract-verifier-backend:
//   beginCell().storeUint(version, 8).storeStringTail(jsonURL).endCell()
struct SourceContent {
    version: uint8
    url: StringTail
}
