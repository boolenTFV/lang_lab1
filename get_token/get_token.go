package get_token
const OUT, EPS, SYM, Q, PQ, OR, REG, CAT, HASH, LP, RP = -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
func GetToken(char byte) int {
	switch(char) {
		case '*': return Q
		case '+': return PQ
		case '|': return OR
		case '(': return LP
		case ')': return RP
		default: return SYM;
	}
}