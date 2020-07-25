package main

import (
	"errors"
)

type pollItem struct {
	level           int
	question        func(user *storageUser, c *ViberCallback) string
	possibleAnswers []string
	validateAnswer  func(string) error
	persistAnswer   func(string, *storageUser) error
}

type poll struct {
	items map[int]*pollItem
	size  int
}

func (p *poll) add(item *pollItem) {
	item.level = p.size
	p.items[item.level] = item
	p.size++
}

func (p *poll) getLevel(level int) *pollItem {
	return p.items[level]
}

func (p *poll) isFinishedFor(u *storageUser) bool {
	return u.Level >= p.size
}

func generateOurPoll() poll {
	ret := poll{
		items: map[int]*pollItem{},
	}

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "Вы гражданин Республики Беларусь?"
		},

		possibleAnswers: []string{
			"1. Да",
			"2. Нет",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["isBelarus"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш возраст"
		},
		possibleAnswers: []string{
			"1. младше 18",
			"2. 18-24",
			"3. 25-34",
			"4. 35-44",
			"5. 45-55",
			"6. старше 55",
		},
		validateAnswer: func(answer string) error {
			if answer == "1. младше 18" {
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

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "К какому типу относится ваш населенный пункт?"
		},
		possibleAnswers: []string{
			"1. Областной центр",
			"2. Город или городской поселок",
			"3. Агрогородок / Село / Деревня",
			"4. Проживаю за пределами РБ",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["residence_location_type"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "За кого Вы планируете проголосовать?"
		},
		possibleAnswers: []string{
			"1. Дмитриев",
			"2. Канопацкая",
			"3. Лукашенко",
			"4. Тихановская",
			"5. Черечень",
			"6. Против всех",
			"7. Затрудняюсь ответить",
			"8. Не пойду голосовать",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Candidate = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш пол"
		},
		possibleAnswers: []string{
			"1. Мужской",
			"2. Женский",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["gender"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "Выберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск"
		},
		possibleAnswers: []string{
			"1. Брестская",
			"2. Витебская",
			"3. Гомельская",
			"4. Гродненская",
			"5. Минская",
			"6. Могилевская",
			"7. Минск",
			"8. Проживаю за пределами РБ",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["residence_location"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "Ваш уровень образования?"
		},
		possibleAnswers: []string{
			"1. Базовое / Среднее общее (школа)",
			"2. Профессионально-техническое",
			"3. Среднее специальное",
			"4. Высшее",
			"5. Другое",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["education_level"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *storageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, общий совокупный месячный доход вашей семьи (включая пенсии, стипендии, надбавки и прочее)"
		},
		possibleAnswers: []string{
			"1. До 500 бел. руб.",
			"2. 500 - 1000 бел. руб.",
			"3. 1000 - 2000 бел. руб.",
			"4. Выше 2000 бел.руб.",
			"5. Не хочу отвечать на этот вопрос",
		},
		persistAnswer: func(answer string, u *storageUser) error {
			u.Properties["income_total"] = answer
			u.isChanged = true
			return nil
		},
	})

	return ret
}
