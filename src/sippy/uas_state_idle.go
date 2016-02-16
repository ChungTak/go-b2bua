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
package sippy

import (
    "sippy/conf"
    "sippy/headers"
    "sippy/types"
)

type UasStateIdle struct {
    uaStateGeneric
    config  sippy_conf.Config
}

func NewUasStateIdle(ua sippy_types.UA, config sippy_conf.Config) *UasStateIdle {
    return &UasStateIdle{
        uaStateGeneric  : newUaStateGeneric(ua),
        config          : config,
    }
}

func (self *UasStateIdle) OnActivation() {
}

func (self *UasStateIdle) String() string {
    return "Idle(UAS)"
}

func (self *UasStateIdle) RecvRequest(req sippy_types.SipRequest, t sippy_types.ServerTransaction) sippy_types.UaState {
    if req.GetMethod() != "INVITE" {
        //print "wrong request %s in the Trying state" % req.getMethod()
        return nil
    }
    self.ua.SetOrigin("caller")
    //print "INVITE received in the Idle state, going to the Trying state"
    if req.GetCGUID() != nil {
        self.ua.SetCGUID(req.GetCGUID().GetCopy())
    } else if req.GetH323ConfId() != nil {
        self.ua.SetH323ConfId(req.GetH323ConfId().GetCopy())
    } else {
        self.ua.SetCGUID(sippy_header.NewSipCiscoGUID())
    }
    self.ua.SetUasResp(req.GenResponse(100, "Trying", nil, self.ua.GetLocalUA().AsSipServer()))
    self.ua.SetLCSeq(100) // XXX: 100 for debugging so that incorrect CSeq generation will be easily spotted
    if self.ua.GetLContact() == nil {
        self.ua.SetLContact(sippy_header.NewSipContact(self.config))
    }
    self.ua.SetRTarget(req.GetContacts()[0].GetUrl().GetCopy())
    self.ua.UpdateRouting(self.ua.GetUasResp(), /*update_rtarget*/ false, /*reverse_routes*/ false)
    self.ua.SetRAddr0(self.ua.GetRAddr())
    t.SendResponse(self.ua.GetUasResp(), false, nil)
    self.ua.GetUasResp().GetTo().SetTag(self.ua.GetLTag())
    self.ua.SetLUri(sippy_header.NewSipFrom(self.ua.GetUasResp().GetTo().GetUri(), self.config))
    self.ua.SetRUri(sippy_header.NewSipTo(self.ua.GetUasResp().GetFrom().GetUri(), self.config))
    self.ua.SetCallId(self.ua.GetUasResp().GetCallId())
    auth := req.GetSipAuthorization().GetCopy()
    body := req.GetBody()
    self.ua.SetBranch(req.GetVias()[0].GetBranch())
    event := NewCCEventTry(self.ua.GetCallId(), self.ua.GetCGUID(), self.ua.GetRUri().GetUrl().Username, req.GetRURI().Username, body, auth,
      self.ua.GetRUri().GetUri().GetName(), /*rtime*/ req.GetRtime(), /*origin*/ self.ua.GetOrigin())
    event.SetReason(req.GetReason())
    if self.ua.GetExpireTime() != 0 {
        self.ua.SetExMtime(event.GetRtime().Add(self.ua.GetExpireTime()))
    }
    if self.ua.GetNoProgressTime() != 0 && (self.ua.GetExpireTime() == 0 || self.ua.GetNoProgressTime() < self.ua.GetExpireTime()) {
        self.ua.SetNpMtime(event.GetRtime().Add(self.ua.GetNoProgressTime()))
    }
    if self.ua.GetNpMtime() != nil {
        self.ua.StartNoProgressTimer(self.ua.GetNpMtime())
    } else if self.ua.GetExMtime() != nil {
        self.ua.StartExpireTimer(self.ua.GetExMtime())
    }
    if body != nil {
        if self.ua.HasOnRemoteSdpChange() {
            self.ua.OnRemoteSdpChange(body, req, func (x sippy_types.MsgBody) { self.ua.DelayedRemoteSdpUpdate(event, x) })
            self.ua.SetSetupTs(req.GetRtime())
            return NewUasStateTrying(self.ua)
        } else {
            self.ua.SetRSDP(body.GetCopy())
        }
    } else {
        self.ua.SetRSDP(nil)
    }
    self.ua.Enqueue(event)
    self.ua.SetSetupTs(req.GetRtime())
    return NewUasStateTrying(self.ua)
}
