#!/bin/zsh

curl -g 'http://localhost:8080/vocab?query={vocab(id:2){learning_lang,first_lang,pos}}'
