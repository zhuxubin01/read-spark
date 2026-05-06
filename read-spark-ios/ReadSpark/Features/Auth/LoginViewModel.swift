import Combine
import Foundation

class LoginViewModel: ObservableObject {
    @Published var phone = ""
    @Published var code = ""
    @Published var isLoading = false
    @Published var error: String?
    @Published var isLoggedIn = false

    private var cancellables = Set<AnyCancellable>()
    private let repository = AuthRepository.shared

    func login() {
        isLoading = true
        error = nil

        repository.login(phone: phone, code: code)
            .sink(receiveCompletion: { [weak self] completion in
                self?.isLoading = false
                if case .failure(let err) = completion {
                    self?.error = err.localizedDescription
                }
            }, receiveValue: { [weak self] _ in
                self?.isLoggedIn = true
            })
            .store(in: &cancellables)
    }
}
