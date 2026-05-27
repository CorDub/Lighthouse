import { type ValueCheck } from "../types/helpersTypes.ts"
import { UserSchema } from "../schemas/user.ts"
import { type ErrorKey } from "../Errors.tsx";

// final functions

export function checkForErrors(...checks: ValueCheck[]): ErrorKey[] {
  const finalErrorList: ErrorKey[] = []

  for (const check of checks) {
    const res = singleCheckForError(check)
    finalErrorList.push(...res)
  }

  return finalErrorList
}

function singleCheckForError(check: ValueCheck): ErrorKey[] {
  const checkFunctionsTable = {
    "email": checkEmail,
    "password": checkPassword
  }

  return checkFunctionsTable[check[0]](check[1])
}



// brick functions

function checkEmail(value: string): ErrorKey[] {
  let errorList: ErrorKey[] = []

  const present = checkRequired(value)
  if (!present) {
    errorList.push("emailRequired")
    return errorList
  }

  const emailValid = checkEmailValid(value)
  if (!emailValid) {
    errorList.push("emailValid")
  }

  const notTooLong = checkMax(value, 254)
  if (!notTooLong) {
    errorList.push("emailTooLong")
  }

  return errorList
}

function checkPassword(value: string): ErrorKey[] {
  let errorList: ErrorKey[] = []

  const present = checkRequired(value)
  if (!present) {
    errorList.push("passwordRequired")
    return errorList
  }

  const notTooShort = checkMin(value, 8)
  if (!notTooShort) {
    errorList.push("passwordTooShort")
  }

  const notTooLong = checkMax(value, 72)
  if (!notTooLong) {
    errorList.push("passwordTooLong")
  }

  return errorList
}



// building block functions

function checkRequired(value: string): boolean {
  if (value.trim().length === 0) {
    return false
  }
  return true
}

function checkEmailValid(value: string): boolean {
  const validityCheck = UserSchema.unwrap().shape.email.safeParse(value)
  if (!validityCheck.success) {
    return false
  }
  return true
}

function checkMax(value: string, max: number): boolean {
  if (value.length > max) {
    return false
  }
  return true
}

function checkMin(value: string, min: number): boolean {
  if (value.length < min) {
    return false
  }
  return true
}