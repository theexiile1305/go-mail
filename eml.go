// SPDX-FileCopyrightText: 2022-2023 The go-mail Authors
//
// SPDX-License-Identifier: MIT

package mail

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	nm "net/mail"
	"os"
	"strings"
)

// EMLToMsg will open an parse a .eml file at a provided file path and return a
// pre-filled Msg pointer
func EMLToMsg(fp string) (*Msg, error) {
	m := &Msg{
		addrHeader:    make(map[AddrHeader][]*nm.Address),
		genHeader:     make(map[Header][]string),
		preformHeader: make(map[Header]string),
		mimever:       MIME10,
	}

	pm, mbbuf, err := readEML(fp)
	if err != nil || pm == nil {
		return m, fmt.Errorf("failed to parse EML file: %w", err)
	}

	if err := parseEMLHeaders(&pm.Header, m); err != nil {
		return m, fmt.Errorf("failed to parse EML headers: %w", err)
	}
	if err := parseEMLBodyParts(pm, mbbuf, m); err != nil {
		return m, fmt.Errorf("failed to parse EML body parts: %w", err)
	}

	return m, nil
}

// readEML opens an EML file and uses net/mail to parse the header and body
func readEML(fp string) (*nm.Message, *bytes.Buffer, error) {
	fh, err := os.Open(fp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open EML file: %w", err)
	}
	defer func() {
		_ = fh.Close()
	}()
	pm, err := nm.ReadMessage(fh)
	if err != nil {
		return pm, nil, fmt.Errorf("failed to parse EML: %w", err)
	}

	buf := bytes.Buffer{}
	if _, err = buf.ReadFrom(pm.Body); err != nil {
		return nil, nil, err
	}
	return pm, &buf, nil
}

// parseEMLHeaders will check the EML headers for the most common headers and set the
// according settings in the Msg
func parseEMLHeaders(mh *nm.Header, m *Msg) error {
	commonHeaders := []Header{
		HeaderContentType, HeaderImportance, HeaderInReplyTo, HeaderListUnsubscribe,
		HeaderListUnsubscribePost, HeaderMessageID, HeaderMIMEVersion, HeaderOrganization,
		HeaderPrecedence, HeaderPriority, HeaderReferences, HeaderSubject, HeaderUserAgent,
		HeaderXMailer, HeaderXMSMailPriority, HeaderXPriority,
	}

	// Extract address headers
	if v := mh.Get(HeaderFrom.String()); v != "" {
		if err := m.From(v); err != nil {
			return fmt.Errorf(`failed to parse "From:" header: %w`, err)
		}
	}
	ahl := map[AddrHeader]func(...string) error{
		HeaderTo:  m.To,
		HeaderCc:  m.Cc,
		HeaderBcc: m.Bcc,
	}
	for h, f := range ahl {
		if v := mh.Get(h.String()); v != "" {
			var als []string
			pal, err := nm.ParseAddressList(v)
			if err != nil {
				return fmt.Errorf(`failed to parse address list: %w`, err)
			}
			for _, a := range pal {
				als = append(als, a.String())
			}
			if err := f(als...); err != nil {
				return fmt.Errorf(`failed to parse "To:" header: %w`, err)
			}
		}
	}

	// Extract date from message
	d, err := mh.Date()
	if err != nil {
		switch {
		case errors.Is(err, nm.ErrHeaderNotPresent):
			m.SetDate()
		default:
			return fmt.Errorf("failed to parse EML date: %w", err)
		}
	}
	if err == nil {
		m.SetDateWithValue(d)
	}

	// Extract common headers
	for _, h := range commonHeaders {
		if v := mh.Get(h.String()); v != "" {
			m.SetGenHeader(h, v)
		}
	}

	return nil
}

// parseEMLBodyParts ...
func parseEMLBodyParts(pm *nm.Message, mbbuf *bytes.Buffer, m *Msg) error {
	// Extract the transfer encoding of the body
	mt, par, err := mime.ParseMediaType(pm.Header.Get(HeaderContentType.String()))
	if err != nil {
		return fmt.Errorf("failed to extract content type: %w", err)
	}
	if v, ok := par["charset"]; ok {
		m.SetCharset(Charset(v))
	}

	cte := pm.Header.Get(HeaderContentTransferEnc.String())
	switch strings.ToLower(mt) {
	case TypeTextPlain.String():
		if cte == NoEncoding.String() {
			m.SetEncoding(NoEncoding)
			m.SetBodyString(TypeTextPlain, mbbuf.String())
			break
		}
		if cte == EncodingQP.String() {
			m.SetEncoding(EncodingQP)
			qpr := quotedprintable.NewReader(mbbuf)
			qpbuf := bytes.Buffer{}
			if _, err := qpbuf.ReadFrom(qpr); err != nil {
				return fmt.Errorf("failed to read quoted-printable body: %w", err)
			}
			m.SetBodyString(TypeTextPlain, qpbuf.String())
			break
		}
		if cte == EncodingB64.String() {
			m.SetEncoding(EncodingB64)
			b64d := base64.NewDecoder(base64.StdEncoding, mbbuf)
			b64buf := bytes.Buffer{}
			if _, err := b64buf.ReadFrom(b64d); err != nil {
				return fmt.Errorf("failed to read base64 body: %w", err)
			}
			m.SetBodyString(TypeTextPlain, b64buf.String())
			break
		}
	default:
	}
	return nil
}
