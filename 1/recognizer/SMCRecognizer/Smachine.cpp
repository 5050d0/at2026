//
// Created by kirill on 3/12/26.
//
#include "Smachine.h"

#include <utility>

Smachine::Smachine() : _fsm(*this) {
    _fsm.enterStartState();
};

Smachine::~Smachine() = default;

void Smachine::SetVartype(std::string str) {
    vartype = std::move(str);
}

void Smachine::reset() {
    success = false;
    foundeq = false;
    foundsign = false;
    vartype = "";
    lvarname = "";
    rvar1name = "";
    rvar2name = "";
}

void Smachine::SetFoundEQ(bool cond) {
    foundeq = cond;
}

void Smachine::SetLvarname(std::string str) {
    lvarname = std::move(str);
}

void Smachine::SetRvar1name(std::string str) {
    rvar1name = std::move(str);
}

void Smachine::SetRvar2name(std::string str) {
    rvar2name = std::move(str);
}

void Smachine::SetLvarname(char ch) {
    lvarname = std::string{ch};
}

void Smachine::SetRvar1name(char ch) {
    rvar1name = std::string{ch};
}

void Smachine::SetRvar2name(char ch) {
    rvar2name = std::string{ch};
}

void Smachine::SetSuccess(bool a) {
    success = a;
}

void Smachine::SetFoundSign(bool b) {
    foundsign = b;
}

bool Smachine::run(const std::string &row) {
    reset();
    _fsm.Reset();

    for (const auto &it: row) {
        _fsm.getchar(it);
    }
    return success;
}
