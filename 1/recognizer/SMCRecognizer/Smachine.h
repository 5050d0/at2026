#pragma once

#include "recognizer_sm.h"

class Smachine {
    recognizerContext _fsm;

public:
    bool success = false;
    bool foundeq = false;
    bool foundsign = false;
    std::string vartype;
    std::string lvarname;
    std::string rvar1name;
    std::string rvar2name;


    Smachine();

    ~Smachine();

    void SetVartype(std::string);

    void reset();

    void SetFoundEQ(bool cond);

    void SetLvarname(std::string);

    void SetRvar1name(std::string);

    void SetRvar2name(std::string);

    void SetLvarname(char);

    void SetRvar1name(char);

    void SetRvar2name(char);

    void SetSuccess(bool);

    void SetFoundSign(bool);

    bool run(const std::string &row);
};


