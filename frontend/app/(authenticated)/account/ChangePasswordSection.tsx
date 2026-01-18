export default function ChangePasswordSection() {
  return (
    <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
      <div>
        <h2 className="font-semibold text-base">Change password</h2>
        <p className="mt-1 text-sm/6 text-gray-400">Update your password associated with your account.</p>
      </div>

      <form className="md:col-span-2">
        <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
          <div className="col-span-full">
            <label htmlFor="current-password" className="block text-sm/6 font-medium text-base">
              Current password
            </label>
            <div className="mt-2">
              <input
                id="current-password"
                name="current_password"
                type="password"
                autoComplete="current-password"
                className="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6"
              />
            </div>
          </div>

          <div className="col-span-full">
            <label htmlFor="new-password" className="block text-sm/6 font-medium text-base">
              New password
            </label>
            <div className="mt-2">
              <input
                id="new-password"
                name="new_password"
                type="password"
                autoComplete="new-password"
                className="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6"
              />
            </div>
          </div>

          <div className="col-span-full">
            <label htmlFor="confirm-password" className="block text-sm/6 font-medium text-base">
              Confirm password
            </label>
            <div className="mt-2">
              <input
                id="confirm-password"
                name="confirm_password"
                type="password"
                autoComplete="new-password"
                className="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:border-accent sm:text-sm/6"
              />
            </div>
          </div>
        </div>

        <div className="mt-8 flex">
          <button
            type="submit"
            className="rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight hover:text-inverted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
          >
            Save
          </button>
        </div>
      </form>
    </div>
  )
}