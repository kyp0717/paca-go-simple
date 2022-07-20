package main

import (
  "fmt"
  )

type Job1 struct {
  paca PacaClient
  features Features
}


func NewJob(ac PacaClient, f Features) Job1 {
  j := Job1{
    paca: ac,
    features: f,
  }
  return j
}

func (j Job1) GetData() Job2 {
  j2 := Job2{}
  snapshot, err := j.paca.data.GetSnapshot(j.features.stock)
  // data 
  if err != nil {
    j2.Metrics = snapshot.MinuteBar.Close
  } else {
    fmt.Println("error")
  }
  return j2
} 



