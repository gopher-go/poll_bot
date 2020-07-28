package main

import (
	"errors"
)

type pollItem struct {
	id              string
	level           int
	question        func(user *storageUser, c *ViberCallback) string
	possibleAnswers []string
	validateAnswer  func(string) error
	persistAnswer   func(string, *storageUser) error
}

type poll struct {
	items map[int]pollItem
	size  int
}

func (p *poll) add(item pollItem) {
	item.level = p.size
	p.items[item.level] = item
	p.size++
}

func (p *poll) getLevel(level int) pollItem {
	return p.items[level]
}

func (p *poll) isFinishedFor(u *storageUser) bool {
	return u.Level >= p.size
}

func generateOurPoll() poll {
	ret := poll{
		items: map[int]pollItem{},
	}

	ret.add(pollItem{
		id: "isBelarus",
		question: func(user *storageUser, c *ViberCallback) string {
			return "Вы гражданин Республики Беларусь?"
		},

		possibleAnswers: []string{
			"Да",
			"Нет",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["isBelarus"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(pollItem{
		id: "age",
		question: func(user *storageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш возраст"
		},
		possibleAnswers: []string{
			"младше 18",
			"18-25",
			"26-40",
			"41-55",
			"55-70",
			"старше 70",
		},
		validateAnswer: func(answer string) error {
			if answer == "младше 18" {
				return errors.New("Вам должно быть 18 или больше.")
			}
			return nil
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["age"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(pollItem{
		id: "residence_location_type",
		question: func(user *storageUser, c *ViberCallback) string {
			return "К какому типу относится населенный пункт, в котором вы проживаете?"
		},
		possibleAnswers: []string{
			"Областной центр / Минск",
			"Город или городской поселок",
			"Агрогородок / Село / Деревня",
			"Проживаю за пределами РБ",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["residence_location_type"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(pollItem{
		id: "candidate",
		question: func(user *storageUser, c *ViberCallback) string {
			return "За кого Вы планируете проголосовать?"
		},
		possibleAnswers: []string{
			"Дмитриев", "Тихановская",
			"Лукашенко", "Канопацкая",
			"Черечень", "Против всех",
			"Затрудняюсь ответить", "Не пойду голосовать",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Candidate = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(pollItem{
		id: "gender",
		question: func(user *storageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш пол"
		},
		possibleAnswers: []string{
			"Мужской",
			"Женский",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["gender"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(pollItem{
		id: "residence_location",
		question: func(user *storageUser, c *ViberCallback) string {
			return "Выберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск"
		},
		possibleAnswers: []string{
			"Брестская",
			"Витебская",
			"Гомельская",
			"Гродненская",
			"Минская",
			"Могилевская",
			"Минск",
			"Проживаю за пределами РБ",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["residence_location"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(pollItem{
		id: "education_level",
		question: func(user *storageUser, c *ViberCallback) string {
			return "Ваш уровень образования?"
		},
		possibleAnswers: []string{
			"Базовое / Среднее общее (школа)",
			"Профессионально-техническое",
			"Среднее специальное",
			"Высшее",
			"Другое",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["education_level"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(pollItem{
		id: "income_total",
		question: func(user *storageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, общий совокупный месячный доход вашей семьи (включая пенсии, стипендии, надбавки и прочее)"
		},
		possibleAnswers: []string{
			"До 500 бел. руб.",
			"500 - 1000 бел. руб.",
			"1000 - 2000 бел. руб.",
			"Выше 2000 бел.руб.",
			"Не хочу отвечать на этот вопрос",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["income_total"] = answer
			u.isChanged = true
			return nil
		},
	})

	return ret
}
