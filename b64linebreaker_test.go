// SPDX-FileCopyrightText: 2022-2023 The go-mail Authors
//
// SPDX-License-Identifier: MIT

package mail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

const logoB64 = `PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+PCFE
T0NUWVBFIHN2ZyBQVUJMSUMgIi0vL1czQy8vRFREIFNWRyAxLjEvL0VOIiAiaHR0cDovL3d3dy53
My5vcmcvR3JhcGhpY3MvU1ZHLzEuMS9EVEQvc3ZnMTEuZHRkIj48c3ZnIHdpZHRoPSIxMDAlIiBo
ZWlnaHQ9IjEwMCUiIHZpZXdCb3g9IjAgMCA2MDAgNjAwIiB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJo
dHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3Jn
LzE5OTkveGxpbmsiIHhtbDpzcGFjZT0icHJlc2VydmUiIHhtbG5zOnNlcmlmPSJodHRwOi8vd3d3
LnNlcmlmLmNvbS8iIHN0eWxlPSJmaWxsLXJ1bGU6ZXZlbm9kZDtjbGlwLXJ1bGU6ZXZlbm9kZDtz
dHJva2UtbGluZWNhcDpyb3VuZDtzdHJva2UtbGluZWpvaW46cm91bmQ7c3Ryb2tlLW1pdGVybGlt
aXQ6MS41OyI+PHJlY3QgaWQ9ItCc0L7QvdGC0LDQttC90LDRjy3QvtCx0LvQsNGB0YLRjDEiIHNl
cmlmOmlkPSLQnNC+0L3RgtCw0LbQvdCw0Y8g0L7QsdC70LDRgdGC0YwxIiB4PSIwIiB5PSIwIiB3
aWR0aD0iNjAwIiBoZWlnaHQ9IjYwMCIgc3R5bGU9ImZpbGw6bm9uZTsiLz48Zz48cGF0aCBkPSJN
NTg0LjcyNiw5OC4zOTVjLTAsLTkuNzM4IC03LjkwNywtMTcuNjQ1IC0xNy42NDUsLTE3LjY0NWwt
NTEwLjAyNywwYy05LjczOSwwIC0xNy42NDUsNy45MDcgLTE3LjY0NSwxNy42NDVsLTAsMzE5Ljc5
NGMtMCw5LjczOSA3LjkwNiwxNy42NDYgMTcuNjQ1LDE3LjY0Nmw1MTAuMDI3LC0wYzkuNzM4LC0w
IDE3LjY0NSwtNy45MDcgMTcuNjQ1LC0xNy42NDZsLTAsLTMxOS43OTRaIiBzdHlsZT0iZmlsbDoj
MmMyYzJjO3N0cm9rZTojMDAwO3N0cm9rZS13aWR0aDoyLjFweDsiLz48cGF0aCBkPSJNNTg1LjEy
Niw5OC4zOTVjLTAsLTkuNzM4IC03LjkwNywtMTcuNjQ1IC0xNy42NDUsLTE3LjY0NWwtNDk0Ljgz
OSwwYy05LjczOCwwIC0xNy42NDUsNy45MDcgLTE3LjY0NSwxNy42NDVsMCwzMTkuNzk0YzAsOS43
MzkgNy45MDcsMTcuNjQ2IDE3LjY0NSwxNy42NDZsNDk0LjgzOSwtMGM5LjczOCwtMCAxNy42NDUs
LTcuOTA3IDE3LjY0NSwtMTcuNjQ2bC0wLC0zMTkuNzk0WiIgc3R5bGU9ImZpbGw6IzM3MzczNztz
dHJva2U6IzAwMDtzdHJva2Utd2lkdGg6Mi4wN3B4OyIvPjxwYXRoIGQ9Ik01NjcuMjU5LDExNC40
NzJjLTAsLTguNzYgLTcuMTEyLC0xNS44NzIgLTE1Ljg3MiwtMTUuODcybC00NjIuNjUxLDBjLTgu
NzYsMCAtMTUuODcyLDcuMTEyIC0xNS44NzIsMTUuODcybDAsMjg3LjY0MWMwLDguNzYgNy4xMTIs
MTUuODcxIDE1Ljg3MiwxNS44NzFsNDYyLjY1MSwwYzguNzYsMCAxNS44NzIsLTcuMTExIDE1Ljg3
MiwtMTUuODcxbC0wLC0yODcuNjQxWiIgc3R5bGU9ImZpbGw6I2ZmZjJlYjtzdHJva2U6IzAwMDtz
dHJva2Utd2lkdGg6MS45cHg7Ii8+PHBhdGggZD0iTTIzNi4xNDUsNDM1Ljg3OGwtMjMuNDIsOTEu
MDA4bDEzMC43NjUsLTBsMjIuMzY1LC05MC43NTZsLTEyOS43MSwtMC4yNTJaIiBzdHlsZT0iZmls
bDojMmMyYzJjO3N0cm9rZTojMDAwO3N0cm9rZS13aWR0aDoyLjQ3cHg7Ii8+PHBhdGggZD0iTTI0
NS4zOTcsNDM1Ljg3OGwtMjMuNDIxLDkxLjAwOGwxMzAuNzY2LC0wbDIyLjM2NSwtOTAuNzU2bC0x
MjkuNzEsLTAuMjUyWiIgc3R5bGU9ImZpbGw6IzM3MzczNztzdHJva2U6IzAwMDtzdHJva2Utd2lk
dGg6Mi40N3B4OyIvPjxwYXRoIGQ9Ik00NzIuNTg0LDUzMS4wOTNjLTAsLTIuODI3IC0yLjI5NSwt
NS4xMjEgLTUuMTIxLC01LjEyMWwtMjkyLjU0MSwtMGMtMi44MjYsLTAgLTUuMTIxLDIuMjk0IC01
LjEyMSw1LjEyMWwwLDEwLjI0M2MwLDIuODI2IDIuMjk1LDUuMTIxIDUuMTIxLDUuMTIxbDI5Mi41
NDEsLTBjMi44MjYsLTAgNS4xMjEsLTIuMjk1IDUuMTIxLC01LjEyMWwtMCwtMTAuMjQzWiIgc3R5
bGU9ImZpbGw6IzJjMmMyYztzdHJva2U6IzAwMDtzdHJva2Utd2lkdGg6Mi40N3B4OyIvPjxnPjxw
YXRoIGQ9Ik0xOTAuMzc5LDQzMS40NTVjLTAsLTEuOTE2IC0xLjU1NiwtMy40NzIgLTMuNDcyLC0z
LjQ3MmwtNTQuNTYzLDBjLTEuOTE3LDAgLTMuNDcyLDEuNTU2IC0zLjQ3MiwzLjQ3MmwtMCw0Ni4z
MDljLTAsMS45MTYgMS41NTUsMy40NzIgMy40NzIsMy40NzJsNTQuNTYzLC0wYzEuOTE2LC0wIDMu
NDcyLC0xLjU1NiAzLjQ3MiwtMy40NzJsLTAsLTQ2LjMwOVoiIHN0eWxlPSJmaWxsOiNlZWFjMDE7
c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjMuMXB4OyIvPjxnPjxwYXRoIGQ9Ik0xMzkuOTk2LDQ0
Mi44NTNjLTAsMC40MiAtMC4wNzEsMC45MjIgMC4wNjMsMS4zMjNjMC41MywxLjU5MiAyLjc1LDEu
Mzg2IDQuMDk1LDEuMzg2YzIuMzk3LC0wIDQuNjc0LC0wLjEwOCA2LjQ4OSwtMS44MjdjMC40NSwt
MC40MjYgMC45NjYsLTAuOTQgMS4wNzIsLTEuNTc1YzAuMDYxLC0wLjM2NyAwLjEwMiwtMS4yMDgg
LTAuMTI2LC0xLjUxMmMtMC44MDMsLTEuMDcxIC0yLjg0NCwtMC4zNzMgLTMuNjU1LDAuMjUyYy0w
LjcxNCwwLjU1IC0wLjgyMiwxLjgyNiAtMC45NDUsMi42NDZjLTAuNDA2LDIuNzA1IDAuMzg3LDUu
MDMgMy4xNSw1Ljg1OWMzLjQwNSwxLjAyMSA2Ljc4NCwwLjA2IDEwLjA4NiwtMC44MDljMS41OTcs
LTAuNDIxIDMuNjYsLTAuNzcgNS4wMTcsLTEuNzJjMC40NTIsLTAuMzE2IDAuNDQ0LC0wLjI4NCAw
LjkwMSwtMC41NThjMS4yMDgsLTAuNzI1IDMuODY5LC0xLjkyNyAyLjQ1NywtMy41OTFjLTAuMTks
LTAuMjI0IC0xLjEzOSwtMC4yMzcgLTEuMjYsMC4xMjZjLTAuNDIsMS4yNTcgLTEuMzI2LDMuMzAy
IDAuNTY3LDMuNzhjMS4zOTgsMC4zNTMgMi45MTksMC4xODkgNC4zNDcsMC4xODljMC4xOTMsLTAg
MS40ODUsLTAuMjc5IDEuNzAxLC0wLjA2M2MwLjU1NSwwLjU1NSAxLjU0MywwLjgxMSAyLjMzMSwx
LjAwOGMwLjkzMiwwLjIzMyAxLjk0NCwwLjYwNiAyLjg5OSwwLjY5M2MwLjg2OSwwLjA3OSAxLjc5
LC0wLjA0NSAyLjY0NiwwLjEyNiIgc3R5bGU9ImZpbGw6bm9uZTtzdHJva2U6IzAwMDtzdHJva2Ut
d2lkdGg6MS45NnB4OyIvPjxwYXRoIGQ9Ik0xNDIuNzA1LDQ1Mi42ODFjMC4xMTksMC41NTEgMC44
OTMsMC42NjkgMS4zMTgsMC44NDZjMS41OTgsMC42NjYgMy41MTQsMC44NzQgNS4yMzQsMC45ODJj
Mi4zNzgsMC4xNDggNC43MDQsLTAuMTU5IDcuMDU3LC0wLjMxNWMxLjg2MSwtMC4xMjUgMy43NDYs
MC4xMTYgNS42MDcsLTBjMi4wNTEsLTAuMTI5IDQuMzIxLC0wLjM3NyA2LjM2NCwtMC4wNjNjMS4w
NzIsMC4xNjQgMi4wODksMC40OCAzLjE1LDAuNjkzIiBzdHlsZT0iZmlsbDpub25lO3N0cm9rZToj
MDAwO3N0cm9rZS13aWR0aDoxLjk2cHg7Ii8+PHBhdGggZD0iTTE3Ni4zNDksNDU1LjA3NmMxLjEx
MywtMC4wMjEgMi4yMjYsLTAuMDQyIDMuMzQsLTAuMDYzIiBzdHlsZT0iZmlsbDpub25lO3N0cm9r
ZTojMDAwO3N0cm9rZS13aWR0aDoxLjk2cHg7Ii8+PHBhdGggZD0iTTE0NS4xNjIsNDU5Ljg2NGwy
MS42NzQsLTAiIHN0eWxlPSJmaWxsOm5vbmU7c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjEuOTZw
eDsiLz48L2c+PC9nPjxnPjxnPjxwYXRoIGQ9Ik03Ni4yNzMsMTIxLjQ5NGw0ODcuNTU3LC0wLjAx
NSIgc3R5bGU9ImZpbGw6bm9uZTtzdHJva2U6I2Y0ZGRkMDtzdHJva2Utd2lkdGg6NC45M3B4O3N0
cm9rZS1saW5lY2FwOnNxdWFyZTsiLz48cGF0aCBkPSJNMjMwLjk1MSwyMjQuMTc1Yy0wLjgwMiwt
My4zMjcgLTguMDU5LC03Ljk4NCAtMTEuMDA1LC04Ljk0NWMtMzUuODQ2LC0xMS42OTkgLTIwLjc4
MywyNy4zNzIgLTAuNjE4LDI0LjM0MiIgc3R5bGU9ImZpbGw6IzYxZTBlZTsiLz48cGF0aCBkPSJN
MjM0Ljk0NCwyMjMuMjEzYy0wLjM3OSwtMS41NzQgLTEuNTIyLC0zLjQzNiAtMy4yNzYsLTUuMTc0
Yy0zLjA2MSwtMy4wMzMgLTguMDE0LC01LjkyIC0xMC40NDgsLTYuNzE0Yy0xMy4xMzUsLTQuMjg3
IC0yMC42MTgsLTEuOTk0IC0yNC4wNzEsMS44MTljLTMuNjM4LDQuMDE3IC0zLjgzNiwxMC40OTIg
LTAuODI1LDE2LjY1MWMzLjk4MSw4LjE0MyAxMy4zMjEsMTUuMzg2IDIzLjYxNSwxMy44MzljMi4y
NDIsLTAuMzM3IDMuNzg4LC0yLjQzIDMuNDUxLC00LjY3MmMtMC4zMzcsLTIuMjQyIC0yLjQzLC0z
Ljc4OSAtNC42NzIsLTMuNDUyYy02LjY0NywwLjk5OSAtMTIuNDQ0LC00LjA2NSAtMTUuMDE0LC05
LjMyM2MtMC44MDEsLTEuNjM3IC0xLjI5MywtMy4zMDggLTEuMzAyLC00Ljg0N2MtMC4wMDUsLTEu
MDIgMC4xOTUsLTEuOTc0IDAuODM2LC0yLjY4MmMwLjgzMiwtMC45MTggMi4yMjUsLTEuMzU0IDQu
MTMzLC0xLjQ4MWMyLjg1MiwtMC4xOTEgNi41NjcsMC40MTIgMTEuMywxLjk1N2MxLjQ1NiwwLjQ3
NSA0LjE2LDIuMDk2IDYuMjU3LDMuODY5YzAuNjM2LDAuNTM3IDEuMjE1LDEuMDg4IDEuNjU4LDEu
NjM2YzAuMTUyLDAuMTg4IDAuMzMxLDAuMzI3IDAuMzcyLDAuNDk4YzAuNTMxLDIuMjA0IDIuNzUx
LDMuNTYyIDQuOTU1LDMuMDMxYzIuMjA0LC0wLjUzMSAzLjU2MiwtMi43NTIgMy4wMzEsLTQuOTU1
WiIvPjxwYXRoIGQ9Ik0zNDkuMzI2LDEyMy41MjdjMCwwIC0xNS4zNSwtMTMuODYyIC0xOC4xMTQs
MjIuODM0IiBzdHlsZT0iZmlsbDpub25lO3N0cm9rZTojNjFlMGVlO3N0cm9rZS13aWR0aDo2LjE2
cHg7Ii8+PHBhdGggZD0iTTI5OS42NjIsMTIzLjQ4NWMtMCwtMCAxMi42NzEsLTE1LjI2NSAyMC4x
MTIsMjEuODk4IiBzdHlsZT0iZmlsbDpub25lO3N0cm9rZTojNjFlMGVlO3N0cm9rZS13aWR0aDo2
LjE2cHg7Ii8+PHBhdGggZD0iTTMzNi4xMTEsMTA4Ljc3M2MtMCwwIC0xNi44MDgsLTkuMDk5IC0x
MC44ODIsMzIuMTI0IiBzdHlsZT0iZmlsbDpub25lO3N0cm9rZTojNjFlMGVlO3N0cm9rZS13aWR0
aDo2LjE2cHg7Ii8+PHBhdGggZD0iTTQxNy43OTYsMTg0LjQ1M2MzLjI4MiwwLjk3IDYuNTgzLC0x
LjI5IDguODg5LC0zLjM2MWMyOC4wNTUsLTI1LjE5NCAtMTMuMTk4LC0zMS41MzggLTIwLjY1Niwt
MTIuNTU5IiBzdHlsZT0iZmlsbDojNjFlMGVlOyIvPjxwYXRoIGQ9Ik00MTYuNjMyLDE4OC4zOTNj
NC42NTcsMS4zNzUgOS41MjYsLTEuMzA3IDEyLjc5NywtNC4yNDVjMTAuMjc5LC05LjIzIDEyLjA0
NiwtMTYuODQzIDEwLjQ4OCwtMjEuNzE2Yy0xLjY0NSwtNS4xNSAtNy4xNDIsLTguNTQ0IC0xMy45
NzgsLTguOTljLTkuMDExLC0wLjU4NyAtMTkuOTI3LDMuOTA0IC0yMy43MzMsMTMuNTg5Yy0wLjgz
LDIuMTEgMC4yMSw0LjQ5NiAyLjMyLDUuMzI1YzIuMTEsMC44MjkgNC40OTYsLTAuMjExIDUuMzI1
LC0yLjMyMWMyLjQ2LC02LjI1OSA5LjczLC04Ljc3NSAxNS41NTMsLTguMzk1YzEuODA3LDAuMTE3
IDMuNDg5LDAuNTE3IDQuODE3LDEuMjY5YzAuODcxLDAuNDk0IDEuNTg2LDEuMTMgMS44NzEsMi4w
MjNjMC4zNzUsMS4xNzIgMC4wNDksMi41ODggLTAuNzk1LDQuMjk4Yy0xLjI2NiwyLjU2MiAtMy42
NDksNS40NzcgLTcuMzU2LDguODA2Yy0wLjg2OCwwLjc4IC0xLjkwNywxLjYxOCAtMy4wNTQsMi4x
NDdjLTAuNjE2LDAuMjg0IC0xLjI2NSwwLjUyNyAtMS45MjgsMC4zMzFjLTIuMTc0LC0wLjY0MiAt
NC40NjEsMC42MDIgLTUuMTAzLDIuNzc2Yy0wLjY0MiwyLjE3NCAwLjYwMiw0LjQ2MSAyLjc3Niw1
LjEwM1oiLz48cGF0aCBkPSJNNDM2LjAxNCw0MTkuNjQ1YzMuNjAxLC0xMC45MzQgOS4wMDMsLTMx
LjYwOSA1LjY4MiwtNDkuNTYzYy0xNC45NTMsLTgwLjg0NyA3Ljg5NSwtOTguNTQ5IC01LjAzNSwt
MTQ3LjU5M2MtNC40OSwtMTcuMDI5IC0xMi4wNDksLTM2LjI1IC0yMy41NTQsLTQ5Ljc5MmMtMjUu
Mjg1LC0yOS43NjIgLTcwLjczNiwtNDIuOTEzIC0xMDguNzMsLTMyLjc4NWMtMTUuMjEzLDQuMDU1
IC0yOC42NjIsMTEuNzI3IC00MS45MjMsMjAuMDM0Yy00Ni4xMzYsMjguOSAtNTIuMjU5LDgzLjM3
NyAtNDIuNzkxLDEyOC4yNDhjNC4xOTEsMTkuODY0IDIzLjY3Myw0OS44OTcgMjIuNTg2LDcwLjE5
NGMtMS4zODcsMjUuOTEgLTAuNzY1LDQ0LjQyNiAzLjM3Niw2MC44ODQiIHN0eWxlPSJmaWxsOiM2
MWUwZWU7Ii8+PGNsaXBQYXRoIGlkPSJfY2xpcDEiPjxwYXRoIGQ9Ik00MzYuMDE0LDQxOS42NDVj
My42MDEsLTEwLjkzNCA5LjAwMywtMzEuNjA5IDUuNjgyLC00OS41NjNjLTE0Ljk1MywtODAuODQ3
IDcuODk1LC05OC41NDkgLTUuMDM1LC0xNDcuNTkzYy00LjQ5LC0xNy4wMjkgLTEyLjA0OSwtMzYu
MjUgLTIzLjU1NCwtNDkuNzkyYy0yNS4yODUsLTI5Ljc2MiAtNzAuNzM2LC00Mi45MTMgLTEwOC43
MywtMzIuNzg1Yy0xNS4yMTMsNC4wNTUgLTI4LjY2MiwxMS43MjcgLTQxLjkyMywyMC4wMzRjLTQ2
LjEzNiwyOC45IC01Mi4yNTksODMuMzc3IC00Mi43OTEsMTI4LjI0OGM0LjE5MSwxOS44NjQgMjMu
NjczLDQ5Ljg5NyAyMi41ODYsNzAuMTk0Yy0xLjM4NywyNS45MSAtMC43NjUsNDQuNDI2IDMuMzc2
LDYwLjg4NCIvPjwvY2xpcFBhdGg+PGcgY2xpcC1wYXRoPSJ1cmwoI19jbGlwMSkiPjxwYXRoIGQ9
Ik0zMTEuNzc2LDIzNi44MzVjMCwwIC0zLjA2OCwyNC4xOTkgLTMxLjU5NywyNC44MzhjLTI4LjUz
LDAuNjM5IC0zMi41OTEsLTI4LjI0NiAtMzIuNTkxLC0yOC4yNDZsNjQuMTg4LDMuNDA4WiIgc3R5
bGU9ImZpbGw6IzRmY2JlMDsiLz48cGF0aCBkPSJNNDMyLjUwOSwyMTYuMDIzYzAsLTAgLTMuMDY3
LDI0LjE5OCAtMzEuNTk3LDI0LjgzN2MtMjguNTMsMC42MzkgLTMyLjU5MSwtMjguMjQ2IC0zMi41
OTEsLTI4LjI0Nmw2NC4xODgsMy40MDlaIiBzdHlsZT0iZmlsbDojNGZjYmUwOyIvPjxwYXRoIGQ9
Ik0yODYuODQ1LDQxOC45NjFsMTI3Ljg4NiwtMGMwLjM4NiwtMS4yNjcgMC42NjYsLTIuNDMyIDAu
ODcsLTMuNDcxYzMuODY3LC0xOS43NDMgMy43MjksLTU3LjQ2OCAtMi40MTUsLTc3LjA5NWMtMi4x
MzMsLTYuODE2IC01Ljc5MywtMTQuNDkgLTExLjQ3NiwtMTkuODUyYy0xMi40OSwtMTEuNzg1IC02
NC4zOTgsLTkuMTgxIC04My41NjEsLTQuODM1Yy03LjY3MywxLjc0MSAtMTQuNDg3LDQuOTIyIC0y
MS4yMSw4LjM1OGMtMTQuOTExLDcuNjIgLTIwLjI1Myw1OC43MTIgLTEwLjA5NCw5Ni44OTVaIiBz
dHlsZT0iZmlsbDojYjdmOGZmOyIvPjxwYXRoIGQ9Ik0yNDYuNzExLDQxOC45NjFjLTUuMTA4LC0y
MC42MTMgLTYuNDksLTQyLjk3NSAtNS41NjIsLTYwLjg1NWMyLjEzNiwtNDEuMTMzIC0zMS4zMTQs
LTUwLjQ1MyAtMjUuNDYxLC0xMjEuNTEzYzE0LjUzNCwtNDIuNDA5IDE1LjA4OSwtNDAuMTA3IDI3
LjY1MywtNTUuNDQ1YzAsMCAtMTMuMTQ3LDM5Ljk2OCAtMS4xOTMsNzAuNzI4YzExLjk1NCwzMC43
NTkgMzIuMTc0LDczLjQxOCAzMi41OSwxMTcuNDkxYzAuMTg2LDE5Ljc3MSAxLjYzNCwzNi4zMSA0
LjEzLDQ5LjU5NGwtMzIuMTU3LC0wWiIgc3R5bGU9ImZpbGw6IzRmY2JlMDsiLz48L2c+PHBhdGgg
ZD0iTTQzOS45MTUsNDIwLjkzYzMuNzQ4LC0xMS4zOCA5LjI3NiwtMzIuOTA4IDUuODIsLTUxLjU5
NWMtOC44NjUsLTQ3LjkzMyAtNC4yNTksLTczLjQwNSAtMS45MzUsLTk2LjU4YzEuNjEzLC0xNi4w
OTMgMi4xNTYsLTMxLjEyMiAtMy4xNjcsLTUxLjMxM2MtNC42MzgsLTE3LjU5MSAtMTIuNTExLC0z
Ny40MTYgLTI0LjM5NiwtNTEuNDA1Yy0yNi4yNjEsLTMwLjkxMSAtNzMuNDU2LC00NC42MTMgLTEx
Mi45MTcsLTM0LjA5NGMtMTUuNjE0LDQuMTYxIC0yOS40MzcsMTEuOTk2IC00My4wNDcsMjAuNTIy
Yy00Ny43MzUsMjkuOTAyIC01NC40MjYsODYuMTUgLTQ0LjYyOSwxMzIuNTc3YzIuMTM5LDEwLjEz
OCA4LjEzNiwyMi44ODggMTMuNTA0LDM1LjY3NmM0Ljk5NCwxMS45IDkuNTE1LDIzLjgxIDguOTk5
LDMzLjQ1Yy0xLjQxNSwyNi40MzIgLTAuNzMsNDUuMzE3IDMuNDk1LDYyLjEwN2wxLjAwMiwzLjk4
M2w3Ljk2NywtMi4wMDVsLTEuMDAyLC0zLjk4M2MtNC4wNTgsLTE2LjEyNyAtNC42MTgsLTM0LjI3
NCAtMy4yNTgsLTU5LjY2M2MwLjQ2LC04LjYwOSAtMi40NiwtMTguODk1IC02LjU0NSwtMjkuNDU3
Yy01LjY5NCwtMTQuNzI0IC0xMy42NDYsLTMwLjA1OSAtMTYuMTI0LC00MS44MDRjLTkuMTQsLTQz
LjMxNSAtMy41ODQsLTk2LjAyMiA0MC45NTIsLTEyMy45MTljMTIuOTEyLC04LjA4OSAyNS45ODks
LTE1LjU5OCA0MC44MDEsLTE5LjU0N2MzNi41MjYsLTkuNzM1IDgwLjIzNCwyLjg2NSAxMDQuNTQx
LDMxLjQ3NmMxMS4xMjYsMTMuMDk2IDE4LjM3MiwzMS43MTMgMjIuNzEzLDQ4LjE4YzYuMzEsMjMu
OTMyIDMuODIsNDAuMjE5IDEuNjQ3LDYwLjM4M2MtMi4yNTcsMjAuOTQ5IC00LjI1OSw0NS45Mjcg
My4zMjEsODYuOTFjMy4xODYsMTcuMjIyIC0yLjA5LDM3LjA0MyAtNS41NDQsNDcuNTMxbC0xLjI4
NSwzLjkwMmw3LjgwMywyLjU2OWwxLjI4NCwtMy45MDFaIi8+PHBhdGggZD0iTTQyNC41MzksMjQy
LjI0N2MtMS44NjgsMC4xMzQgLTMuMjc2LDEuNzU4IC0zLjE0MiwzLjYyNmMwLjEzMywxLjg2NyAx
Ljc1NywzLjI3NSAzLjYyNSwzLjE0MWMxLjg2NywtMC4xMzMgMy4yNzUsLTEuNzU3IDMuMTQyLC0z
LjYyNWMtMC4xMzQsLTEuODY3IC0xLjc1OCwtMy4yNzUgLTMuNjI1LC0zLjE0MloiIHN0eWxlPSJm
aWxsOiNiN2Y4ZmY7Ii8+PHBhdGggZD0iTTI1Ni45OTcsMjYwLjM5MWMtMS44NjgsMC4xMzMgLTMu
Mjc1LDEuNzU4IC0zLjE0MiwzLjYyNWMwLjEzMywxLjg2OCAxLjc1OCwzLjI3NSAzLjYyNSwzLjE0
MmMxLjg2OCwtMC4xMzMgMy4yNzUsLTEuNzU4IDMuMTQyLC0zLjYyNWMtMC4xMzMsLTEuODY4IC0x
Ljc1OCwtMy4yNzUgLTMuNjI1LC0zLjE0MloiIHN0eWxlPSJmaWxsOiNiN2Y4ZmY7Ii8+PGc+PHBh
dGggZD0iTTM2NC4wOSwyNjkuMjI5Yy0xLjE4Niw0Ljk2NyAtMS40NDMsMTQuMzA2IC00LjMxNSwx
OC40Yy0xLjM5MywxLjk4NSAtNC40OTUsMi4wODYgLTYuNTE2LDIuMTUzYy0xMS42MTIsMC4zODMg
LTExLjQ0NiwtMTEuNTEzIC0xMy4xODEsLTIxLjg1NSIgc3R5bGU9ImZpbGw6I2VkZWRlZDsiLz48
cGF0aCBkPSJNMzYxLjA5NCwyNjguNTEzYy0wLjcxOSwzLjAwOSAtMS4xMDYsNy42MDcgLTEuOTIy
LDExLjY5OGMtMC40NDUsMi4yMjggLTAuOTcxLDQuMjk4IC0xLjkxOSw1LjY0OWMtMC4yNjMsMC4z
NzYgLTAuNzQ5LDAuNDU5IC0xLjIwMywwLjU2NmMtMC45OTUsMC4yMzMgLTIuMDU2LDAuMjUgLTIu
ODkzLDAuMjc3Yy0yLjA1NSwwLjA2OCAtMy42MSwtMC4zMjYgLTQuNzc0LC0xLjE5MmMtMS4xOTMs
LTAuODg3IC0xLjk1NCwtMi4yMTIgLTIuNTQ0LC0zLjc0M2MtMS41NjEsLTQuMDQ3IC0xLjg5LC05
LjM5IC0yLjcyMywtMTQuMzUxYy0wLjI4MSwtMS42NzcgLTEuODcxLC0yLjgxIC0zLjU0OCwtMi41
MjhjLTEuNjc3LDAuMjgxIC0yLjgwOSwxLjg3MSAtMi41MjgsMy41NDhjMC45MDMsNS4zODEgMS4z
NTgsMTEuMTU4IDMuMDUxLDE1LjU0OGMxLjAzNiwyLjY4OCAyLjUyLDQuOTEyIDQuNjE1LDYuNDdj
Mi4xMjUsMS41ODEgNC45MDMsMi41MyA4LjY1NCwyLjQwNmMxLjQ2OCwtMC4wNDggMy40LC0wLjE1
MiA1LjA3OSwtMC43MTRjMS41NiwtMC41MjIgMi45MTksLTEuNDEgMy44NTgsLTIuNzQ5YzEuMzUs
LTEuOTI0IDIuMjg0LC00LjgwOCAyLjkxNywtNy45ODFjMC44LC00LjAxMSAxLjE2OCwtOC41MjIg
MS44NzMsLTExLjQ3M2MwLjM5NSwtMS42NTMgLTAuNjI3LC0zLjMxNyAtMi4yODEsLTMuNzEyYy0x
LjY1NCwtMC4zOTUgLTMuMzE3LDAuNjI3IC0zLjcxMiwyLjI4MVoiLz48cGF0aCBkPSJNMzQ5LjAw
MSwyMzQuMDVjLTkuMzMyLDEuMjQ2IC0xNi4yNjUsNy4wNzUgLTE1LjQ3MywxMy4wMDhjMC43OTIs
NS45MzMgOS4wMTIsOS43MzggMTguMzQ0LDguNDkyYzkuMzMzLC0xLjI0NiAxNi4yNjYsLTcuMDc1
IDE1LjQ3NCwtMTMuMDA4Yy0wLjc5MiwtNS45MzMgLTkuMDEyLC05LjczOCAtMTguMzQ1LC04LjQ5
MloiLz48cGF0aCBkPSJNMzQ4LjkwMSwyNTIuNTVjMC40Niw1LjIwNCAxLjY5NCwxMC4xNDIgMi40
NTMsMTUuMjc3YzAuMjQ4LDEuNjgyIDEuODE1LDIuODQ2IDMuNDk3LDIuNTk4YzEuNjgyLC0wLjI0
OCAyLjg0NiwtMS44MTYgMi41OTgsLTMuNDk4Yy0wLjc0MSwtNS4wMTUgLTEuOTYxLC05LjgzNyAt
Mi40MTEsLTE0LjkyYy0wLjE1LC0xLjY5NCAtMS42NDcsLTIuOTQ3IC0zLjM0LC0yLjc5N2MtMS42
OTQsMC4xNSAtMi45NDcsMS42NDcgLTIuNzk3LDMuMzRaIi8+PHBhdGggZD0iTTM3NC4yNiwyNTku
ODE0Yy04LjAyMiw1LjM1IC0xNS45MDEsNy44MjUgLTIzLjU5NCw3LjI3NGMtNy42OTQsLTAuNTUx
IC0xNS4xNDgsLTQuMTI0IC0yMi4zNDIsLTEwLjU4OWMtMS4yNjUsLTEuMTM2IC0zLjIxNCwtMS4w
MzIgLTQuMzUxLDAuMjMzYy0xLjEzNiwxLjI2NCAtMS4wMzIsMy4yMTQgMC4yMzMsNC4zNWM4LjM1
Niw3LjUwOCAxNy4wODQsMTEuNTExIDI2LjAyLDEyLjE1MWM4LjkzNiwwLjY0IDE4LjEzNSwtMi4w
NzggMjcuNDUzLC04LjI5NGMxLjQxNCwtMC45NDMgMS43OTcsLTIuODU4IDAuODUzLC00LjI3MmMt
MC45NDMsLTEuNDE1IC0yLjg1OCwtMS43OTcgLTQuMjcyLC0wLjg1M1oiLz48cGF0aCBkPSJNMzc4
LjYwNSwxNTIuMTU2Yy0zLjY0NywtMC44MjYgLTYuOTU2LDAuMDM4IC03LjM4NSwxLjkyOWMtMC40
MjgsMS44OTEgMi4xODUsNC4wOTcgNS44MzIsNC45MjRjMy42NDcsMC44MjcgNi45NTYsLTAuMDM3
IDcuMzg1LC0xLjkyOGMwLjQyOCwtMS44OTEgLTIuMTg1LC00LjA5OCAtNS44MzIsLTQuOTI1WiIg
c3R5bGU9ImZpbGw6I2I3ZjhmZjsiLz48cGF0aCBkPSJNMjQ5LjUsMTc3LjI3Yy0yLjE5NSwzLjAy
OCAtMi43MDIsNi40MSAtMS4xMzIsNy41NDhjMS41NywxLjEzOCA0LjYyNiwtMC4zOTYgNi44MjEs
LTMuNDI0YzIuMTk1LC0zLjAyOCAyLjcwMiwtNi40MSAxLjEzMiwtNy41NDhjLTEuNTcsLTEuMTM4
IC00LjYyNiwwLjM5NiAtNi44MjEsMy40MjRaIiBzdHlsZT0iZmlsbDojYjdmOGZmOyIvPjxwYXRo
IGQ9Ik00MDMuNjQ4LDE2Mi44MjVjLTE5LjIxMywtMS4zNSAtMzUuOTA3LDEzLjE1MSAtMzcuMjU3
LDMyLjM2NGMtMS4zNSwxOS4yMTIgMTMuMTUyLDM1LjkwNiAzMi4zNjQsMzcuMjU2YzE5LjIxMywx
LjM1MSAzNS45MDcsLTEzLjE1MSAzNy4yNTcsLTMyLjM2NGMxLjM1LC0xOS4yMTIgLTEzLjE1Miwt
MzUuOTA2IC0zMi4zNjQsLTM3LjI1NloiIHN0eWxlPSJmaWxsOiNmZmY7c3Ryb2tlOiMwMDA7c3Ry
b2tlLXdpZHRoOjYuMTZweDsiLz48cGF0aCBkPSJNMjgxLjY5NiwxODEuNzhjLTE5LjIxMiwtMS4z
NSAtMzUuOTA2LDEzLjE1MiAtMzcuMjU3LDMyLjM2NGMtMS4zNSwxOS4yMTIgMTMuMTUyLDM1Ljkw
NyAzMi4zNjUsMzcuMjU3YzE5LjIxMiwxLjM1IDM1LjkwNiwtMTMuMTUyIDM3LjI1NiwtMzIuMzY0
YzEuMzUsLTE5LjIxMyAtMTMuMTUyLC0zNS45MDcgLTMyLjM2NCwtMzcuMjU3WiIgc3R5bGU9ImZp
bGw6I2ZmZjtzdHJva2U6IzAwMDtzdHJva2Utd2lkdGg6Ni4xNnB4OyIvPjxwYXRoIGQ9Ik0zOTgu
NDk5LDE5Mi44ODFjLTYuODI0LC0wLjQ3OSAtMTIuNzU0LDQuNjcyIC0xMy4yMzMsMTEuNDk2Yy0w
LjQ4LDYuODI0IDQuNjcxLDEyLjc1MyAxMS40OTUsMTMuMjMzYzYuODI0LDAuNDc5IDEyLjc1NCwt
NC42NzIgMTMuMjMzLC0xMS40OTZjMC40OCwtNi44MjQgLTQuNjcxLC0xMi43NTMgLTExLjQ5NSwt
MTMuMjMzWiIvPjxwYXRoIGQ9Ik0yNzguMTQyLDIxMC4yMjdjLTYuODI0LC0wLjQ3OSAtMTIuNzU0
LDQuNjcxIC0xMy4yMzMsMTEuNDk2Yy0wLjQ4LDYuODI0IDQuNjcxLDEyLjc1MyAxMS40OTUsMTMu
MjMzYzYuODI0LDAuNDc5IDEyLjc1NCwtNC42NzIgMTMuMjMzLC0xMS40OTZjMC40OCwtNi44MjQg
LTQuNjcxLC0xMi43NTMgLTExLjQ5NSwtMTMuMjMzWiIvPjxwYXRoIGQ9Ik0zOTkuODkyLDE5Ny45
NjFjLTEuOTE2LC0wLjEzNSAtMy41ODEsMS4zMTEgLTMuNzE2LDMuMjI3Yy0wLjEzNCwxLjkxNiAx
LjMxMiwzLjU4MSAzLjIyOCwzLjcxNmMxLjkxNiwwLjEzNCAzLjU4MSwtMS4zMTIgMy43MTUsLTMu
MjI4YzAuMTM1LC0xLjkxNiAtMS4zMTEsLTMuNTgxIC0zLjIyNywtMy43MTVaIiBzdHlsZT0iZmls
bDojZmZmOyIvPjxwYXRoIGQ9Ik0yNzkuNTM1LDIxNS4zMDZjLTEuOTE2LC0wLjEzNCAtMy41ODEs
MS4zMTIgLTMuNzE1LDMuMjI4Yy0wLjEzNSwxLjkxNiAxLjMxMSwzLjU4MSAzLjIyNywzLjcxNmMx
LjkxNiwwLjEzNCAzLjU4MSwtMS4zMTIgMy43MTYsLTMuMjI4YzAuMTM0LC0xLjkxNiAtMS4zMTIs
LTMuNTgxIC0zLjIyOCwtMy43MTZaIiBzdHlsZT0iZmlsbDojZmZmOyIvPjwvZz48L2c+PHBhdGgg
ZD0iTTUxMi4yOTEsMzMzLjcyNmMwLC0wIC0yNy40OTYsMTMuNDAyIC00My4yNDUsMy43OTJjLTE1
Ljc1LC05LjYxIC0xMC44LC0yNC43ODggMy4xMDIsLTIyLjUyYzEzLjkwMywyLjI2OCAyNy4xNjgs
My4xOTUgMzUuOTQzLC0zLjkzN2M4Ljc3NiwtNy4xMzIgMjEuNjMzLDE0LjM3NiA0LjIsMjIuNjY1
WiIgc3R5bGU9ImZpbGw6I2Q2YzI5NTtzdHJva2U6IzAwMDtzdHJva2Utd2lkdGg6Ni4xNnB4OyIv
PjwvZz48Zz48cGF0aCBkPSJNNDQxLjI2OCwyOTEuMzQ5Yy0wLDAgMjguMDgsNDcuNDg4IC0xMjIu
NTkxLDExNC4wNGMtMTEuNjI5LC0yLjAxNyAtMTQuMzMsLTguNTE1IC0xNC4zMywtOC41MTVjMCwt
MCA0My42MTcsLTEyLjU4NiA4My44MTIsLTQ0Ljg3NmM0My41OTcsLTM1LjAyNCA1My4xMDksLTYw
LjY0OSA1My4xMDksLTYwLjY0OVoiIHN0eWxlPSJmaWxsOiM3YjQ3MjI7c3Ryb2tlOiMwMDA7c3Ry
b2tlLXdpZHRoOjYuNjhweDsiLz48cGF0aCBkPSJNMjMzLjk1NywzNzIuMTM4YzMuNTYsLTUuNTA2
IDEwLjcyOSwtNy4zOTIgMTYuNTM4LC00LjM1MWMxNy43NjcsOS4yMzUgNTMuNjE0LDI4LjAzMiA2
Ny4xNTQsMzYuMzc2YzUuOTcxLDMuNjggLTYuNDc3LDE2LjkxMyAtMTkuMjgzLDM3LjkyN2MtMy4z
NjgsNS41NjcgLTEwLjQzNiw3LjYyOCAtMTYuMjY5LDQuNzQ2Yy0yMS42NTYsLTEwLjcxNiAtNDEu
MDQ3LC0yMC45NTggLTU0LjIzNCwtMjguMDczYy00LjU5MiwtMi41MDEgLTcuOTM0LC02LjgwMiAt
OS4yMjMsLTExLjg2OWMtMS4yODgsLTUuMDY4IC0wLjQwNywtMTAuNDQzIDIuNDMzLC0xNC44MzNj
NC4zMzUsLTYuNzA4IDkuMTYyLC0xNC4xNyAxMi44ODQsLTE5LjkyM1oiIHN0eWxlPSJmaWxsOiM3
YjQ3MjI7c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjQuOTVweDsiLz48cGF0aCBkPSJNMjM5Ljg1
LDM2MS4xNzFsLTI1Ljg0OCwyMy4yOTZsNzkuMTcsNDUuNjA1bDI1LjE5NiwtMjQuMjk1bC03OC41
MTgsLTQ0LjYwNloiIHN0eWxlPSJmaWxsOiNhZDZjNGE7c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRo
OjQuOTVweDsiLz48cGF0aCBkPSJNMjM4LjA4MSwzNzcuNzA0YzIuNTI2LDAuNTMgNC4xNDcsMy4w
MTEgMy42MTgsNS41MzdjLTAuNTMsMi41MjYgLTMuMDExLDQuMTQ3IC01LjUzOCwzLjYxOGMtMi41
MjYsLTAuNTMgLTQuMTQ3LC0zLjAxMSAtMy42MTcsLTUuNTM4YzAuNTMsLTIuNTI2IDMuMDExLC00
LjE0NyA1LjUzNywtMy42MTdaIiBzdHlsZT0iZmlsbDojZmZiMjVjO3N0cm9rZTojMDAwO3N0cm9r
ZS13aWR0aDo0Ljk1cHg7Ii8+PHBhdGggZD0iTTI5MC4zMyw0MDcuOWMyLjUyNywwLjUzIDQuMTQ3
LDMuMDExIDMuNjE4LDUuNTM3Yy0wLjUzLDIuNTI3IC0zLjAxMSw0LjE0NyAtNS41MzcsMy42MThj
LTIuNTI3LC0wLjUzIC00LjE0OCwtMy4wMTEgLTMuNjE4LC01LjUzN2MwLjUzLC0yLjUyNyAzLjAx
MSwtNC4xNDggNS41MzcsLTMuNjE4WiIgc3R5bGU9ImZpbGw6I2ZmYjI1YztzdHJva2U6IzAwMDtz
dHJva2Utd2lkdGg6NC45NXB4OyIvPjwvZz48Zz48cGF0aCBkPSJNMzA1LjMyMywzMzYuMjM1YzAs
MCAyNC4yMDQsMTEyLjU3OCAzNS4zNTcsMTE0LjM1NWMxMS4xNTMsMS43NzYgMTgyLjk5OCwtNTcu
MTYyIDE4Mi45OTgsLTU3LjE2MmMwLC0wIC0yNC4yNywtMTA1Ljc3MyAtMzQuNDU2LC0xMTQuNjAy
Yy0xMC4xODYsLTguODMgLTE4My44OTksNTcuNDA5IC0xODMuODk5LDU3LjQwOVoiIHN0eWxlPSJm
aWxsOiNmZmYyZWI7c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjQuNjhweDsiLz48cGF0aCBkPSJN
NDEzLjU2MywzNjcuNjU2bC0xMDguMjQsLTMxLjQyMWwxMDguMjQsMzEuNDIxWiIgc3R5bGU9ImZp
bGw6bm9uZTtzdHJva2U6IzAwMDtzdHJva2Utd2lkdGg6NS40NnB4OyIvPjxwYXRoIGQ9Ik00ODgu
NTY2LDI3OS4xNTZsLTc1LjAwMyw4OC41IiBzdHlsZT0iZmlsbDpub25lO3N0cm9rZTojMDAwO3N0
cm9rZS13aWR0aDo1LjQ2cHg7Ii8+PHBhdGggZD0iTTMzOS4xNjcsNDQ5LjEyN2w2Mi40MDEsLTgz
LjUwNSIgc3R5bGU9ImZpbGw6bm9uZTtzdHJva2U6IzAwMDtzdHJva2Utd2lkdGg6NS40NnB4OyIv
PjxwYXRoIGQ9Ik00MjMuMjcsMzU1LjgzNWw5OS4xOTQsMzcuMDAxIiBzdHlsZT0iZmlsbDpub25l
O3N0cm9rZTojMDAwO3N0cm9rZS13aWR0aDo1LjQ2cHg7Ii8+PC9nPjxwYXRoIGQ9Ik0zMDQuNDE2
LDM3NC43ODVjLTAsMCAtMzAuNDYsLTIuOCAtMzguOTQzLC0xOS4xODVjLTguNDgyLC0xNi4zODQg
My42MjIsLTI2Ljc5NCAxNC4zMzIsLTE3LjY0NmMxMC43MTEsOS4xNDkgMjEuNTcyLDE2LjgyMSAz
Mi43NzQsMTUuMjc0YzExLjIwMSwtMS41NDcgMTEuMDQsMjMuNTExIC04LjE2MywyMS41NTdaIiBz
dHlsZT0iZmlsbDojZDZjMjk1O3N0cm9rZTojMDAwO3N0cm9rZS13aWR0aDo2LjE2cHg7Ii8+PHBh
dGggZD0iTTEwOS45NTEsMzExLjcxMmw2OC4wNzYsLTAiIHN0eWxlPSJmaWxsOm5vbmU7c3Ryb2tl
OiMwMDA7c3Ryb2tlLXdpZHRoOjQuOTVweDsiLz48cGF0aCBkPSJNMTM5LjMwNSwzNTEuNjgzbDY4
LjA3NiwwIiBzdHlsZT0iZmlsbDpub25lO3N0cm9rZTojMDAwO3N0cm9rZS13aWR0aDo0Ljk1cHg7
Ii8+PHBhdGggZD0iTTEzOS42MTcsMzMxLjY5OGwzOC4wOTgsLTAiIHN0eWxlPSJmaWxsOm5vbmU7
c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjQuOTVweDsiLz48cGF0aCBkPSJNODUuMjAxLDIxMS4y
NTNjLTAsLTIuMjEzIC0xLjc5NywtNC4wMSAtNC4wMSwtNC4wMWwtNTMuNDg3LC0wYy0yLjIxNCwt
MCAtNC4wMTEsMS43OTcgLTQuMDExLDQuMDFsMCw1OC44OTRjMCwyLjIxMyAxLjc5Nyw0LjAxIDQu
MDExLDQuMDFsNTMuNDg3LDBjMi4yMTMsMCA0LjAxLC0xLjc5NyA0LjAxLC00LjAxbC0wLC01OC44
OTRaIiBzdHlsZT0iZmlsbDojZWVhYzAxO3N0cm9rZTojMDAwO3N0cm9rZS13aWR0aDozLjQzcHg7
Ii8+PHBhdGggZD0iTTU2NC42NDMsNDI2LjU5MWMwLC0yLjIxMyAtMS43OTYsLTQuMDEgLTQuMDEs
LTQuMDFsLTUzLjQ4NywtMGMtMi4yMTMsLTAgLTQuMDEsMS43OTcgLTQuMDEsNC4wMWwwLDU4Ljg5
NGMwLDIuMjEzIDEuNzk3LDQuMDEgNC4wMSw0LjAxbDUzLjQ4NywwYzIuMjE0LDAgNC4wMSwtMS43
OTcgNC4wMSwtNC4wMWwwLC01OC44OTRaIiBzdHlsZT0iZmlsbDojZWVhYzAxO3N0cm9rZTojMDAw
O3N0cm9rZS13aWR0aDozLjQzcHg7Ii8+PGNpcmNsZSBjeD0iNTE3LjAxNCIgY3k9IjEwOC43NzEi
IHI9IjQuMzUyIiBzdHlsZT0iZmlsbDojZDNkMmIzOyIvPjxjaXJjbGUgY3g9IjUzMS4zMiIgY3k9
IjEwOC43NzEiIHI9IjQuMzUyIiBzdHlsZT0iZmlsbDojZjRiMzAzOyIvPjxjaXJjbGUgY3g9IjU0
NC43OTIiIGN5PSIxMDguNzcxIiByPSI0LjM1MiIgc3R5bGU9ImZpbGw6I2Y5NjEzZDsiLz48cGF0
aCBkPSJNMTA1LjI1Nyw0ODQuNDNjLTAsMCAyMC4zNzQsLTI1LjIyNSA0MS43MTksLTcuNzYxYzIx
LjM0NCwxNy40NjMgLTEyLjYxMyw1MC40NSAtNDAuNzQ5LDQyLjY4OSIgc3R5bGU9ImZpbGw6bm9u
ZTtzdHJva2U6IzAwMDtzdHJva2Utd2lkdGg6OS45cHg7Ii8+PHBhdGggZD0iTTExMy4wMTksNDYw
LjE0MWMtMCwtNC4yNjUgLTMuNDYzLC03LjcyNyAtNy43MjgsLTcuNzI3bC03OC42NTUsLTBjLTQu
MjY1LC0wIC03LjcyNywzLjQ2MiAtNy43MjcsNy43MjdsLTAsNzguNjU1Yy0wLDQuMjY1IDMuNDYy
LDcuNzI4IDcuNzI3LDcuNzI4bDc4LjY1NSwtMGM0LjI2NSwtMCA3LjcyOCwtMy40NjMgNy43Mjgs
LTcuNzI4bC0wLC03OC42NTVaIiBzdHlsZT0iZmlsbDojNjFlMGVlO3N0cm9rZTojMDAwO3N0cm9r
ZS13aWR0aDo0LjU1cHg7c3Ryb2tlLWxpbmVjYXA6c3F1YXJlOyIvPjxwYXRoIGQ9Ik03Ny4yNzcs
Mzk1LjQwMmMtMi4yMzQsOC43NTUgLTYuNDk2LDE1LjM2NiAtMTEuNjQ5LDE3LjcxNGMxLjk5NSwz
LjA1NyAzLjE2Nyw2Ljc4MSAzLjE2NywxMC43OTdjMCwxMC4zODcgLTcuODQsMTguODIgLTE3LjQ5
OCwxOC44MmMtOS42NTcsMCAtMTcuNDk4LC04LjQzMyAtMTcuNDk4LC0xOC44MmMwLC05LjczMyA2
Ljg4NCwtMTcuNzUxIDE1LjY5NywtMTguNzIyYy0zLjM5OSwtNS45NiAtNS41MjQsLTE0LjQyMyAt
NS41MjQsLTIzLjgwNWMwLC0xOC4wMjQgNy44NDEsLTMyLjY1NyAxNy40OTksLTMyLjY1N2M2Ljc5
MiwtMCAxMi42ODUsNy4yMzggMTUuNTg0LDE3LjgwMmMxLjYyOSwtMC41OTIgMy4zODcsLTAuOTE1
IDUuMjIsLTAuOTE1YzguNDUsMCAxNS4zMTEsNi44NjEgMTUuMzExLDE1LjMxMWMtMCw4LjQ1IC02
Ljg2MSwxNS4zMTEgLTE1LjMxMSwxNS4zMTFjLTEuNzUsMCAtMy40MzIsLTAuMjk0IC00Ljk5OCwt
MC44MzZaIiBzdHlsZT0iZmlsbDojZmZmO2ZpbGwtb3BhY2l0eTowLjU7Ii8+PHBhdGggZD0iTTMz
LjU4MywyMzMuOTNjMi41MjksLTAuMzM3IDQuNTEsLTEuMTc4IDYuNzE2LC0yLjQwM2MwLjk1Nywt
MC41MzIgMS45NjksLTAuOTA0IDIuODk4LC0xLjQ4NWMwLjc4OSwtMC40OTMgMS4zNDMsLTAuOTc2
IDEuMzQzLC0xLjk3OWMwLC0wLjIzNiAwLjIyNCwtMC42MzIgMCwtMC43MDdjLTAuOTUyLC0wLjMx
NyAtMi4zNjQsMC40NjUgLTMuMDQsMC45OWMtMC42MjUsMC40ODYgLTEuNDk0LDAuNTc5IC0yLjEy
LDEuMDZjLTEuMDM2LDAuNzk2IC0xLjY4MiwyLjQ2MyAtMS43NjcsMy43NDZjLTAuMDYsMC44OTYg
LTAuMzU3LDIuODQgMC41NjUsMy4zOTNjMC4zNTcsMC4yMTQgMC45MjksMC4xNDIgMS4zNDMsMC4x
NDJjMS4wOCwtMCAyLjEwNCwwLjAxMiAzLjE4MSwtMC4xNDJjMi4wODIsLTAuMjk3IDQuMTY1LC0x
LjQ1NSA2LjAwOSwtMi40MDNjMi4wMTQsLTEuMDM2IDQuMDEyLC0xLjkyMSA1LjU4NCwtMy42MDVj
MC42NDIsLTAuNjg4IDEuMDEsLTEuMjg2IDEuMjcyLC0yLjE5MWMwLjA1OSwtMC4yMDMgLTAuMDYs
LTAuNzE4IDAuMDcxLC0wLjg0OWMwLjEwOCwtMC4xMDggMC4wOTksMC4yMjYgMC4wNzEsMC4yODNj
LTAuMTg1LDAuMzY3IC0wLjQ2MiwwLjc0IC0wLjYzNywxLjEzMWMtMC4zMjQsMC43MjkgLTAuMzUz
LDEuNzU5IC0wLjM1MywyLjU0NWMwLDAuMTMyIC0wLjA3LDAuNTcxIDAsMC42MzZjMC4wNTUsMC4w
NTEgMC4zNjMsMC4xOTMgMC40MjQsMC4yMTJjMi4wMDIsMC42MzMgNC4zNDIsLTAuNTMyIDYuMDA5
LC0xLjQ4NGMwLjY0MSwtMC4zNjcgMS4zNjMsLTAuODUgMi4wNiwtMS4xMTJjMC4wMjYsLTAuMDEg
MC42NzgsLTAuMTk3IDAuNjk2LC0wLjE2MWMwLjI4MiwwLjU2MyAwLjExMSwxLjUxMSAwLjIxMiwy
LjEyMWMwLjI4NSwxLjcwNyAyLjM2NywxLjYxNSAzLjY3NiwxLjQ4NGMwLjIwNywtMC4wMiAwLjM2
MywtMC4yMTIgMC41NjYsLTAuMjEyIiBzdHlsZT0iZmlsbDpub25lO3N0cm9rZTojMDAwO3N0cm9r
ZS13aWR0aDoyLjc1cHg7Ii8+PHBhdGggZD0iTTM2LjQ4MiwyNDIuMzQyYzEuMzQ2LC0wIDIuNjgy
LC0wLjE1NCA0LjAyOSwtMC4yMTJjNC4zMzIsLTAuMTg5IDguNjcxLC0wLjE5NiAxMy4wMDYsLTAu
MjgzYzYuNTI1LC0wLjEzMSAxMy4wNTYsMC4wNzEgMTkuNTgxLDAuMDcxIiBzdHlsZT0iZmlsbDpu
b25lO3N0cm9rZTojMDAwO3N0cm9rZS13aWR0aDoyLjc1cHg7Ii8+PHBhdGggZD0iTTUxNi4wOSw0
NDIuMDMxYy0wLDAuOTM2IC0wLjUzMiwxLjg5MiAtMC4xMjYsMi44MzljMC41MjYsMS4yMjggMi4z
NjMsMS4yNzUgMy40NjksMC45NDZjMC42MjgsLTAuMTg3IDEuMzE2LC0xLjE3NiAxLjcwMywtMS42
NGMxLjAzOCwtMS4yNDcgMi40MzUsLTIuMzc1IDIuNzEyLC00LjAzN2MwLjA0MSwtMC4yNDggMC4x
OTQsLTAuOTM2IDAuMDYzLC0xLjE5OGMtMC43NjMsLTEuNTI1IC0zLjA1OSwtMS40NzggLTQuMTY4
LC0wLjM2OWMtMC4xNjQsMC4xNjQgLTAuNDQ2LDAuMjI3IC0wLjU2MywwLjQzMmMtMC45OTQsMS43
NDUgLTAuNDM5LDQuMzQxIDAuNDQyLDUuOTkyYzAuMzA0LDAuNTY5IDAuODI5LDEuMDQ4IDEuMzI0
LDEuNDUxYzAuMzM5LDAuMjc1IDAuNzc1LDAuNzE2IDEuMTk5LDAuOTI4YzMuMjYyLDEuNjMxIDcu
Njg5LC0wLjEzMyAxMC41OTYsLTEuNzQ4YzEuOTUyLC0xLjA4NSA0LjEsLTIuMTUyIDQuMSwtNC42
NjhjLTAsLTAuMzI1IDAuMDM3LC0wLjY5NCAtMC4wNjMsLTEuMDA5Yy0wLjM2MSwtMS4xMzMgLTIu
MzA3LC0xLjA0MyAtMy4xNTQsLTAuNjk0Yy0xLjE3OSwwLjQ4NiAtMi4zNzMsMS4zOTUgLTMuMDks
Mi40NmMtMS4xNzQsMS43NDIgLTEuMjA3LDQuMjgzIDAuNTY3LDUuNjE0YzEuMjU2LDAuOTQxIDIu
OTU4LDAuOTIzIDQuNDQ3LDAuNzY2YzIuMjQ3LC0wLjIzNyA0LjI1MywtMS4zMTMgNi4wMjMsLTIu
NjU5YzAuODc1LC0wLjY2NCAyLjA5NCwtMS4zMTYgMi42NDksLTIuMzMzYzAuMTczLC0wLjMxNiAw
LjIxNSwtMC43NDUgMC4zNzksLTEuMDczYzAuMDcxLC0wLjE0MyAtMC4xNTgsMC4yODUgLTAuMTg5
LDAuNDQyYy0wLjA1NiwwLjI3NiAtMC4wOTQsMC44NzggLTAsMS4xMzVjMC4zNzYsMS4wMjggMS4y
MzMsMS41MDkgMi4yNywxLjcwM2MxLjM0NiwwLjI1MiAzLjAwOSwtMC4wMSA0LjM1MiwtMC4xODlj
MC45MzksLTAuMTI1IDEuODkyLC0wLjUxNyAxLjg5MiwwLjc1NyIgc3R5bGU9ImZpbGw6bm9uZTtz
dHJva2U6IzAwMDtzdHJva2Utd2lkdGg6Mi40NXB4OyIvPjxwYXRoIGQ9Ik01MTQuMjU5LDQ2MC4z
NTljMC4yNjYsMC4xMzQgMC4wOCwxLjM1NCAwLjU2NywxLjM4OGMwLjQ0MiwwLjAzMSAwLjg4NSww
LjA0IDEuMzI4LDAuMDVjMC41MTcsMC4wMTIgMS4xMDQsMC4wOTggMS42MzMsLTAuMDM0YzEuNjg4
LC0wLjQyMiAzLjc5NiwtMS41OTggNC40ODIsLTMuMjk2YzAuMTg4LC0wLjQ2NyAwLjI0NiwtMS4w
OTQgMC4zMTUsLTEuNTc3YzAuMDc1LC0wLjUyNCAwLjM3OCwtMS43MTggLTAuMDYzLC0yLjE0NGMt
MC4wODgsLTAuMDg1IC0wLjYwNCwtMC4wNjMgLTAuNjk0LC0wLjA2M2MtMC42OTUsLTAgLTEuNzk2
LDAuMzkyIC0yLjIwNywxLjAwOWMtMC4yNzMsMC40MDkgLTAuMTg5LDAuODUxIC0wLjE4OSwxLjMy
NGMtMCwwLjYyIC0wLjAzMywxLjIyNCAwLjEyNiwxLjgzYzAuMzYyLDEuMzc5IDIuMTExLDEuOTQx
IDMuMzQzLDIuMDE4YzIuNjk2LDAuMTY4IDUuNDQ3LC0wLjI0NSA4LjEzNiwtMGMyLjIwNywwLjIg
NC40NTQsMC43MTcgNi42ODYsMC41MDRjMS41OCwtMC4xNSAzLjE1NCwtMC41MjUgNC43MywtMC42
M2MzLjM5NCwtMC4yMjYgNi43OTEsLTAuNTA1IDEwLjIxOCwtMC41MDUiIHN0eWxlPSJmaWxsOm5v
bmU7c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjIuNDVweDsiLz48cGF0aCBkPSJNNTE5LjA1NCw0
NzEuNjhsMCwwLjk0NmMwLDAuMTU1IC0wLDAuMjUxIDAuMDYzLDAuMzc5bDAsMC4xMjZjMC42MDIs
LTIuNDA3IDMuNzc1LC0yLjc4NCA1LjgxLC0yLjc4NGMwLjE1LC0wIDAuMjMyLDAuNDggMC4yNDUs
MC41NzZjMC4wNTUsMC4zODUgLTAuMTkyLDEuNDUgMC4zNzksMS42NGMwLjk4OSwwLjMzIDIuNDA2
LC0wLjMxOCAzLjM0MywtMC41NjhjMC4yMzUsLTAuMDYyIDEuMDA5LC0wLjQ0MSAxLjE5OCwtMC4y
NTJjMC45NTQsMC45NTQgMi40MTIsMC41NjIgMy43MjEsMC42MzFjMy4zMzIsMC4xNzUgNi42Mjks
LTAuMTg5IDkuOTY2LC0wLjE4OWMxLjA1OCwtMCAyLjE2MiwwLjA4NyAzLjIxNywtMGMwLjA2OCwt
MC4wMDYgMS4wMDksLTAuMDI3IDEuMDA5LC0wLjI1MyIgc3R5bGU9ImZpbGw6bm9uZTtzdHJva2U6
IzAwMDtzdHJva2Utd2lkdGg6Mi40NXB4OyIvPjwvZz48Zz48cGF0aCBkPSJNMzc0LjU3OCwxMTQu
NDU1Yy0wLjYxMywtNC4yMTYgLTQuNTMzLC03LjE0MiAtOC43NDksLTYuNTNsLTk5Ljk0NSwxNC41
MTVjLTQuMjE2LDAuNjEzIC03LjE0Miw0LjUzMyAtNi41Myw4Ljc0OWwyLjIxOSwxNS4yNzhjMC42
MTIsNC4yMTYgNC41MzIsNy4xNDIgOC43NDksNi41M2w5OS45NDUsLTE0LjUxNWM0LjIxNiwtMC42
MTIgNy4xNDIsLTQuNTMzIDYuNTI5LC04Ljc0OWwtMi4yMTgsLTE1LjI3OFoiIHN0eWxlPSJmaWxs
OiM5ZGFjY2Q7c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjUuMjFweDsiLz48cGF0aCBkPSJNMzEz
Ljg2OSwxMTcuNDE0YzAsMCA2Ni4xMDIsOS44NzggNzcuNTgyLC0yMy43MzNjMTEuNDgsLTMzLjYx
MSAtMTIuOTg3LC01MS44NzMgLTQxLjYzNywtMzkuOTIxYy0yOC42NSwxMS45NTIgLTU0Ljg1LDQ2
LjE0MiAtNTQuODUsNDYuMTQyYzAsMCAtMzAuMDA1LC04LjEwOCAtNDEuMzE0LDUuMjIxYy0xMS4z
MDksMTMuMzMgLTE3LjE2NywyNi42NDYgLTIuNDM1LDMxLjUxOWMxNC43MzIsNC44NzMgNjIuNjU0
LC0xOS4yMjggNjIuNjU0LC0xOS4yMjhaIiBzdHlsZT0iZmlsbDojNzM1ZDkzO3N0cm9rZTojMDAw
O3N0cm9rZS13aWR0aDo1LjkzcHg7Ii8+PHBhdGggZD0iTTM3NC4xNDgsMTM2Ljg3NmMtMCwwIDM4
LjQ5MywtOS4zMzggMzYuNTYzLC0zLjU0OWMtMS45Myw1Ljc5IC0zNy42MzgsMjYuOTAxIC0xMTQu
ODMxLDE1Ljk2NmM3MC4zNDUsLTEwLjczMSA3OC4yNjgsLTEyLjQxNyA3OC4yNjgsLTEyLjQxN1oi
IHN0eWxlPSJmaWxsOiM1MDNhNzE7c3Ryb2tlOiMwMDA7c3Ryb2tlLXdpZHRoOjQuOTVweDsiLz48
cGF0aCBkPSJNMzU3LjM0NiwxMjQuOTA2bC01NC4wMjEsOS42MTJsLTAsNS4zOTlsNTQuMDIxLC05
LjYxMmwwLC01LjM5OVoiIHN0eWxlPSJmaWxsOiNmZmY7Ii8+PGNpcmNsZSBjeD0iMzAwLjE1IiBj
eT0iMTM3LjE5OCIgcj0iNC42NzciIHN0eWxlPSJmaWxsOiNmZmIyNWM7c3Ryb2tlOiMwMDA7c3Ry
b2tlLXdpZHRoOjQuOTVweDsiLz48Y2lyY2xlIGN4PSIzNjEuNDEiIGN5PSIxMjcuNTgxIiByPSI0
LjI4MiIgc3R5bGU9ImZpbGw6I2ZmYjI1YztzdHJva2U6IzAwMDtzdHJva2Utd2lkdGg6NC45NXB4
OyIvPjwvZz48L3N2Zz4=
`

// TestBase64LineBreaker tests the Write and Close methods of the Base64LineBreaker
func TestBase64LineBreaker(t *testing.T) {
	l, err := os.Open("assets/gopher2.svg")
	if err != nil {
		t.Errorf("failed to open gopher logo asset: %s", err)
		return
	}
	defer func() { _ = l.Close() }()

	var wbuf bytes.Buffer
	lb := Base64LineBreaker{out: &wbuf}
	bw := base64.NewEncoder(base64.StdEncoding, &lb)
	if _, err := io.Copy(bw, l); err != nil {
		t.Errorf("failed to write logo asset to line breaker: %s", err)
		return
	}
	if err := bw.Close(); err != nil {
		t.Errorf("failed to close b64 encoder: %s", err)
	}
	if err := lb.Close(); err != nil {
		t.Errorf("failed to close line breaker: %s", err)
	}
	ob := removeNewLines([]byte(logoB64))
	nb := removeNewLines(wbuf.Bytes())
	if string(ob) != string(nb) {
		t.Errorf("generated line breaker output differs from original data")
	}
}

// TestBase64LineBreakerFailures tests the cases in which the Base64LineBreaker would fail
func TestBase64LineBreakerFailures(t *testing.T) {
	stt := []byte("short")
	ltt := []byte(logoB64)

	// No output writer defined
	lb := Base64LineBreaker{}
	if _, err := lb.Write(stt); err == nil {
		t.Errorf("writing to Base64LineBreaker with no output io.Writer was supposed to failed, but didn't")
	}
	if err := lb.Close(); err != nil {
		t.Errorf("failed to close Base64LineBreaker: %s", err)
	}

	// Closed output writer
	wbuf := errorWriter{}
	fb := Base64LineBreaker{out: wbuf}
	if _, err := fb.Write(ltt); err == nil {
		t.Errorf("writing to Base64LineBreaker with errorWriter was supposed to failed, but didn't")
	}
	if err := fb.Close(); err != nil {
		t.Errorf("failed to close Base64LineBreaker: %s", err)
	}
}

func TestBase64LineBreaker_WriteAndClose(t *testing.T) {
	tests := []struct {
		name          string
		data          []byte
		expectedWrite string
	}{
		{
			name:          "Write data within MaxBodyLength",
			data:          []byte("testdata"),
			expectedWrite: "testdata",
		},
		{
			name:          "Write data exceeds MaxBodyLength",
			data:          []byte("verylongtestdata"),
			expectedWrite: "verylongtest",
		},
	}

	var mockErr = errors.New("mock write error")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			blr := &Base64LineBreaker{out: &buf}
			mw := &mockWriter{writeError: mockErr}
			blr.out = mw

			_, err := blr.Write(tt.data)
			if err != nil && !errors.Is(err, mockErr) {
				t.Errorf("Unexpected error while writing: %v", err)
				return
			}
			err = blr.Close()
			if err != nil && !errors.Is(err, mockErr) {
				t.Errorf("Unexpected error while closing: %v", err)
				return
			}
		})
	}
}

// removeNewLines removes any newline characters from the given data
func removeNewLines(data []byte) []byte {
	result := make([]byte, len(data))
	n := 0

	for _, b := range data {
		if b == '\r' || b == '\n' {
			continue
		}
		result[n] = b
		n++
	}

	return result[0:n]
}

type errorWriter struct{}

func (e errorWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("supposed to always fail")
}

func (e errorWriter) Close() error {
	return fmt.Errorf("supposed to always fail")
}

// MockWriter is a mock implementation of io.Writer used for testing.
type mockWriter struct {
	writeError error
}

// Write writes the data into a buffer.
func (w *mockWriter) Write(p []byte) (n int, err error) {
	return 0, w.writeError
}

func FuzzBase64LineBreaker_Write(f *testing.F) {
	f.Add([]byte("abc"))
	f.Add([]byte("def"))
	f.Add([]uint8{0o0, 0o1, 0o2, 30, 255})
	buf := bytes.Buffer{}
	bw := bufio.NewWriter(&buf)
	f.Fuzz(func(t *testing.T, data []byte) {
		b := &Base64LineBreaker{out: bw}
		if _, err := b.Write(data); err != nil {
			t.Errorf("failed to write to B64LineBreaker: %s", err)
		}
		if err := b.Close(); err != nil {
			t.Errorf("failed to close B64LineBreaker: %s", err)
		}
	})
}
