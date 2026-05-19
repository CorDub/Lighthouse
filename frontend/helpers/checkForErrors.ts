import { type ValueCheck } from "./helpersTypes.ts"
import { UserSchema } from "../src/schemas/user.ts"


// final functions

export function checkForErrors(...checks: ValueCheck[]): string[] {
  console.log("check for Errors")
  console.log("checks", checks)
  const finalErrorList: string[] = []

  for (const check of checks) {
    const res = singleCheckForError(check)
    finalErrorList.push(...res)
  }

  return finalErrorList
}

function singleCheckForError(check: ValueCheck): string[] {
  const checkFunctionsTable = {
    "email": checkEmail,
    "password": checkPassword
  }

  return checkFunctionsTable[check[0]](check[1])
}



// brick functions

function checkEmail(value: string): string[] {
  let errorList: string[] = []

  const present = checkRequired(value)
  if (!present) {
    errorList.push("Email is required")
    return errorList
  }

  const emailValid = checkEmailValid(value)
  if (!emailValid) {
    errorList.push("Email is not valid")
  }

  const notTooLong = checkMax(value, 254)
  if (!notTooLong) {
    errorList.push("Email is too long - max 254 charcaters")
  }

  return errorList
}

function checkPassword(value: string): string[] {
  let errorList: string[] = []

  const present = checkRequired(value)
  if (!present) {
    errorList.push("Password is required")
    return errorList
  }

  const notTooShort = checkMin(value, 8)
  if (!notTooShort) {
    errorList.push("Password is too short - min 8 characters")
  }

  const notTooLong = checkMax(value, 72)
  if (!notTooLong) {
    errorList.push("Password is too long - max 72 characters")
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