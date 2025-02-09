// Copyright 2021 The go-frogeum Authors
// This file is part of the go-frogeum library.
//
// The go-frogeum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-frogeum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-frogeum library. If not, see <http://www.gnu.org/licenses/>.

package snap

import (
	"time"

	"github.com/frogeum/go-frogeum/p2p/tracker"
)

// requestTracker is a singleton tracker for request times.
var requestTracker = tracker.New(ProtocolName, time.Minute)
