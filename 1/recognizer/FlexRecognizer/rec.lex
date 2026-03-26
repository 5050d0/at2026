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
    UNKNOWN,
    TOKEN_EOF
};

struct Token {
    TokenType type;
    std::string text;
};

%}

%option noyywrap
%option c++

VARTYPE   ("int"|"short"|"long")
VARNAME   ([a-zA-Z][a-zA-Z0-9]{0,15})
NUMBER    ([0-9]+)
SIGN      ([%/*])
SPACE     [[:space:]]+
EQUALS    "="
SEMICOLON ";"

%%

{VARTYPE}   { return static_cast<int>(TokenType::VARTYPE); }
{VARNAME}   { return static_cast<int>(TokenType::VARNAME); }
{NUMBER}    { return static_cast<int>(TokenType::NUMBER); }
{SIGN}      { return static_cast<int>(TokenType::SIGN); }
{EQUALS}    { return static_cast<int>(TokenType::EQUALS); }
{SEMICOLON} { return static_cast<int>(TokenType::SEMICOLON); }
{SPACE}     {}
<<EOF>>     { return static_cast<int>(TokenType::TOKEN_EOF); }
.           { return static_cast<int>(TokenType::UNKNOWN); }

%%


std::optional<std::array<std::string, 4>> extract(const std::string& line) {
    std::array<std::string, 4> data;
    std::istringstream in_stream(line);

    yyFlexLexer lexer(&in_stream);

    TokenType type = static_cast<TokenType>(lexer.yylex());
    if (type != TokenType::VARTYPE) {
            return std::nullopt;
    }
    data[0] = lexer.YYText();

    if (static_cast<TokenType>(lexer.yylex()) != TokenType::VARNAME){
                return std::nullopt;
    }
    data[1] = lexer.YYText();
    if (static_cast<TokenType>(lexer.yylex()) != TokenType::EQUALS){
        return std::nullopt;
    }

    type = static_cast<TokenType>(lexer.yylex());
    if (type != TokenType::VARNAME && type != TokenType::NUMBER){
        return std::nullopt;
    }
    if (type == TokenType::VARNAME){
        data[2] = lexer.YYText();
    }
    type = static_cast<TokenType>(lexer.yylex());
    if (type == TokenType::SEMICOLON) {
         if (static_cast<TokenType>(lexer.yylex()) != TokenType::TOKEN_EOF){
                return std::nullopt;
            }
         return data;
    }
    if (type!=TokenType::SIGN) {
        return std::nullopt;
    }
    type = static_cast<TokenType>(lexer.yylex());
    if (type != TokenType::VARNAME && type != TokenType::NUMBER){
        return std::nullopt;
    }
    if (type==TokenType::VARNAME) {
        data[3] = lexer.YYText();
    }
    type = static_cast<TokenType>(lexer.yylex());
    if (type == TokenType::SEMICOLON) {
         if (static_cast<TokenType>(lexer.yylex()) != TokenType::TOKEN_EOF){
                return std::nullopt;
            }
         return data;
    }

    return std::nullopt;
}