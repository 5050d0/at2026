//
// Created by kirill on 3/11/26.
//

#include "SmcRecognizer.h"

std::optional<RecResult> SmcRecognizer::Recognize(std::string row) {
    bool const success = machine.run(row);
    if (success) {
        RecResult result{
            .vartype = machine.vartype, .lvar = machine.lvarname, .rvar1 = machine.rvar1name, .rvar2 = machine.rvar2name
        };
        return result;
    }
    return std::nullopt;
}

void SmcRecognizer::reset() {
    machine.reset();
}
