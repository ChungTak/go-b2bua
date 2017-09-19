// Copyright (c) 2003-2005 Maxim Sobolev. All rights reserved.
// Copyright (c) 2006-2015 Sippy Software, Inc. All rights reserved.
// Copyright (c) 2015 Andrii Pylypenko. All rights reserved.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation and/or
// other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
// ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package sippy_header

import (
    "sippy/conf"
    "sippy/utils"
)

type sipAddressWithTag struct {
    address *sipAddress
}

func ParseSipAddressWithTag(body string, config sippy_conf.Config) (*sipAddressWithTag, error) {
    addr, err := ParseSipAddress(body, true /* relaxedparser */, config)
    if err != nil { return nil, err }
    return &sipAddressWithTag{ address : addr }, nil
}

func NewSipAddressWithTag(address *sipAddress, config sippy_conf.Config) *sipAddressWithTag {
    if address == nil {
        address = NewSipAddress("Anonymous", NewSipURL("" /* username */,
                                    config.GetMyAddress(),
                                    config.GetMyPort(),
                                    false))
    }
    return &sipAddressWithTag{ address : address }
}

func (self *sipAddressWithTag) GenTag() {
    self.address.SetParam("tag", sippy_utils.GenTag())
}

func (self *sipAddressWithTag) GetTag() string {
    return self.address.GetParam("tag")
}

func (self *sipAddressWithTag) SetTag(value string) {
    if value != "" {
        self.address.SetParam("tag", value)
    }
}

func (self *sipAddressWithTag) getCopy() *sipAddressWithTag {
    return &sipAddressWithTag{
        address : self.address.GetCopy(),
    }
}

func (self *sipAddressWithTag) GetUrl() *SipURL {
    return self.address.url
}

func (self *sipAddressWithTag) GetUri() *sipAddress {
    return self.address
}