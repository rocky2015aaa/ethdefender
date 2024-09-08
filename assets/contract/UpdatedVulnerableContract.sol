// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract UpdatedVulnerableContract is Initializable, OwnableUpgradeable, PausableUpgradeable, UUPSUpgradeable, ReentrancyGuardUpgradeable {
    mapping(address => uint256) public balances;

    // Initializer function replaces the constructor
    function initialize() public initializer {
        __Ownable_init(msg.sender);
        __Pausable_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
    }

    // Receive function to accept ETH and update balances
    receive() external payable whenNotPaused {
        balances[msg.sender] += msg.value;
        emit Received(msg.sender, msg.value);
    }

    // Withdraw function to allow users to withdraw their ETH
    function withdraw() external whenNotPaused nonReentrant {
        uint256 amount = balances[msg.sender];
        require(amount > 0, "Insufficient balance");

        // Update balance before sending Ether
        balances[msg.sender] = 0;

        // Safe transfer of Ether using transfer
        payable(msg.sender).transfer(amount);
    }

    // Function to pause the contract, callable only by the owner
    function pause() external onlyOwner {
        _pause();
    }

    // Function to unpause the contract, callable only by the owner
    function unpause() external onlyOwner {
        _unpause();
    }

    // Function to get the current owner(Added function)
    function getOwner() external view returns (address) {
        return owner();
    }

    // Authorization for contract upgrades
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    // Event for receiving ETH
    event Received(address indexed sender, uint256 amount);
}
