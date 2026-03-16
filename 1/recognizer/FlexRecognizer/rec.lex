%{
#include <array>
#include <optional>
#include <sstream>


struct ParseState {
    std::array<std::string, 4> data;
    int index = 0;
    bool complete = false;
    bool valid = true;
};

ParseState state;
%}

%option noyywrap
%option c++

VARTYPE "int"|"short"|"long"
VARNAME [a-zA-Z][a-zA-Z0-9]{0,15}
NUMBER  [0-9]+
SIGN    [%/*]
SPACE[ \t\n\r]+
EQUALS  "="
SEMICOLON ";"

%%

{VARTYPE}   {
    if (state.index == 0) {
        state.data[0] = YYText();
        state.index++;
    } else {
        state.valid = false;
    }
}

{VARNAME}   {
    if (state.index != 0 && state.index < 4) {
        state.data[state.index] = YYText();
        state.index++;
    } else {
        state.valid = false;
    }
}

{NUMBER}    {
    if (state.index > 1 && state.index < 4) {
        state.index++;
    } else {
        state.valid = false;
    }
}

{SEMICOLON} {
    state.complete = true;
}

{SIGN}      {
    if (state.index != 3) {
        state.valid = false;
    }
}

{EQUALS}    {
    if (state.index != 2) {
        state.valid = false;
    }
}

{SPACE}     {}
.           { state.valid = false; }

%%

std::optional<std::array<std::string, 4>> extract(const std::string& line) {
    state = ParseState();

    std::istringstream input_stream(line);
    yyFlexLexer lexer(&input_stream);

    lexer.yylex();

    if (state.complete && state.valid && state.index >2) {
        return state.data;
    }

    return std::nullopt;
}