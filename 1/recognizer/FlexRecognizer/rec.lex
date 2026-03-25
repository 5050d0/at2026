%{
#include <array>
#include <optional>
#include <sstream>
#include <vector>
#include <string>

enum class TokenType {
    VARTYPE = 1,
    VARNAME,
    NUMBER,
    SIGN,
    EQUALS,
    SEMICOLON,
    UNKNOWN
};

struct Token {
    TokenType type;
    std::string text;
};

static std::vector<Token> current_tokens;
%}

%option noyywrap
%option c++

VARTYPE   ("int"|"short"|"long")
VARNAME   ([a-zA-Z][a-zA-Z0-9]{0,15})
NUMBER    ([0-9]+)
SIGN      ([%/*])
SPACE     ([ \t\n\r]+)
EQUALS    ("=")
SEMICOLON (";")

%%

{VARTYPE}   { current_tokens.push_back({TokenType::VARTYPE, YYText()}); }
{VARNAME}   { current_tokens.push_back({TokenType::VARNAME, YYText()}); }
{NUMBER}    { current_tokens.push_back({TokenType::NUMBER, YYText()}); }
{SIGN}      { current_tokens.push_back({TokenType::SIGN, YYText()}); }
{EQUALS}    { current_tokens.push_back({TokenType::EQUALS, YYText()}); }
{SEMICOLON} { current_tokens.push_back({TokenType::SEMICOLON, YYText()}); }
{SPACE}     {}
.           { current_tokens.push_back({TokenType::UNKNOWN, YYText()}); }

%%

bool isVarOrNum(TokenType type) {
    return type == TokenType::VARNAME || type == TokenType::NUMBER;
}

class StringLexer : public yyFlexLexer {
    const char* input_ptr;
    size_t input_len;

public:
    StringLexer(const std::string& str)
        : input_ptr(str.data()), input_len(str.length()) {}

    int LexerInput(char* buf, int max_size) override {
        int to_copy = std::min(max_size, static_cast<int>(input_len));
        if (to_copy > 0) {
            std::memcpy(buf, input_ptr, to_copy);
            input_ptr += to_copy;
            input_len -= to_copy;
            return to_copy;
        }
        return 0;
    }
};

std::optional<std::array<std::string, 4>> extract(const std::string& line) {
    current_tokens.clear();

    StringLexer lexer(line);

    lexer.yylex();

    std::array<std::string, 4> data;

    if (current_tokens.size() == 5) {
        if (current_tokens[0].type == TokenType::VARTYPE &&
            current_tokens[1].type == TokenType::VARNAME &&
            current_tokens[2].type == TokenType::EQUALS &&
            isVarOrNum(current_tokens[3].type) &&
            current_tokens[4].type == TokenType::SEMICOLON) {

            data[0] = current_tokens[0].text;
            data[1] = current_tokens[1].text;

            data[2] = (current_tokens[3].type == TokenType::VARNAME) ? current_tokens[3].text : "";
            data[3] = "";

            return data;
        }
    }
    else if (current_tokens.size() == 7) {
        if (current_tokens[0].type == TokenType::VARTYPE &&
            current_tokens[1].type == TokenType::VARNAME &&
            current_tokens[2].type == TokenType::EQUALS &&
            isVarOrNum(current_tokens[3].type) &&
            current_tokens[4].type == TokenType::SIGN &&
            isVarOrNum(current_tokens[5].type) &&
            current_tokens[6].type == TokenType::SEMICOLON) {

            data[0] = current_tokens[0].text;
            data[1] = current_tokens[1].text;

            data[2] = (current_tokens[3].type == TokenType::VARNAME) ? current_tokens[3].text : "";
            data[3] = (current_tokens[5].type == TokenType::VARNAME) ? current_tokens[5].text : "";

            return data;
        }
    }
    return std::nullopt;
}