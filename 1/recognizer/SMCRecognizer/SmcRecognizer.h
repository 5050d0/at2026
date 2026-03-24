//
// Created by kirill on 3/11/26.
//

#ifndef LAB1_SMCRECOGNIZER_H
#define LAB1_SMCRECOGNIZER_H

#include "Smachine.h"
#include "recognizer/IRecognizer.h"

class SmcRecognizer : public IRecognizer {
    Smachine machine;

public:
    std::optional<RecResult> Recognize(const std::string &row) override;

    void reset() override;
};


#endif //LAB1_SMCRECOGNIZER_H
