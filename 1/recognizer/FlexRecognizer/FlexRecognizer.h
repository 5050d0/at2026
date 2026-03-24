//
// Created by kirill on 3/12/26.
//

#ifndef LAB1_FLEXRECOGNIZER_H
#define LAB1_FLEXRECOGNIZER_H

#include "recognizer/IRecognizer.h"
#include "lex.yy.h"

class FlexRecognizer : public IRecognizer {
public:
    std::optional<RecResult> Recognize(const std::string &row) override {
        std::optional<std::array<std::string, 4> > ex = extract(row);
        if (ex.has_value()) {
            auto t = ex.value();
            RecResult res{.vartype = t[0], .lvar = t[1], .rvar1 = t[2], .rvar2 = t[3]};
            return res;
        }
        return std::nullopt;
    }

    void reset() override {
    }
};


#endif //LAB1_FLEXRECOGNIZER_H
